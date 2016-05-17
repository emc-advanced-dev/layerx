package task_launcher

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/emc-advanced-dev/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/layerx-core/rpi_messenger"
)

type TaskLauncher struct {
	state *lxstate.State
}

func NewTaskLauncher(state *lxstate.State) *TaskLauncher {
	return &TaskLauncher{
		state: state,
	}
}

func (tl *TaskLauncher) LaunchStagedTasks() error {
	stagingTasks, err := tl.state.StagingTaskPool.GetTasks()
	if err != nil {
		return lxerrors.New("retrieving staging tasks", err)
	}

	nodeTaskMap := make(map[string][]*lxtypes.Task)
	for _, task := range stagingTasks {
		nodeId := task.NodeId
		_, ok := nodeTaskMap[nodeId]
		if !ok {
			nodeTaskMap[nodeId] = []*lxtypes.Task{}
		}
		nodeTaskMap[nodeId] = append(nodeTaskMap[nodeId], task)
	}
	for nodeId, tasksToLaunch := range nodeTaskMap {
		nodeResourcePool, err := tl.state.NodePool.GetNodeResourcePool(nodeId)
		if err != nil {
			return lxerrors.New("finding resource pool for node "+nodeId, err)
		}
		resourcesToUseMap, err := nodeResourcePool.GetResources()
		if err != nil {
			return lxerrors.New("retrieving resource list for node "+nodeId, err)
		}
		resourceRpiMap := make(map[string][]*lxtypes.Resource)
		for _, resource := range resourcesToUseMap {
			resourceRpiMap[resource.RpiName] = append(resourceRpiMap[resource.RpiName], resource)
		}
		logrus.WithFields(logrus.Fields{
			"tasks":            fmt.Sprintf("%v", tasksToLaunch),
			"resources":        fmt.Sprintf("%v", resourcesToUseMap),
			"resources_by_rpi": fmt.Sprintf("%v", resourceRpiMap),
			"node_id":          fmt.Sprintf("%s", nodeId),
			"rpi_urls":         fmt.Sprintf("%s", tl.state.GetRpiUrls()),
		}).Debugf("attempting to launch tasks on rpi")

		for rpiName, resourcesToUse := range resourceRpiMap {
			rpi, err := tl.state.RpiPool.GetRpi(rpiName)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"tasks":     fmt.Sprintf("%v", tasksToLaunch),
					"resources": fmt.Sprintf("%v", resourcesToUseMap),
					"node_id":   fmt.Sprintf("%s", nodeId),
					"rpi":       rpiName,
				}).Errorf("retreiving rpi for name ")
				return lxerrors.New("retreiving rpi for name", err)
			}
			err = rpi_messenger.SendLaunchTasksMessage(rpi.Url, tasksToLaunch, resourcesToUse)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"tasks":     fmt.Sprintf("%v", tasksToLaunch),
					"resources": fmt.Sprintf("%v", resourcesToUseMap),
					"node_id":   fmt.Sprintf("%s", nodeId),
					"rpi_url":   rpi.Url,
				}).Errorf("trying to launch tasks on rpi")
				return lxerrors.New("sending launch task message to rpi", err)
			}
			//flush resources from node
			for resourceId, resource := range resourcesToUseMap {
				err := nodeResourcePool.DeleteResource(resourceId)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"resource": fmt.Sprintf("%v", resource),
						"node_id":  fmt.Sprintf("%s", nodeId),
					}).Errorf("flushing resource " + resourceId + " from node " + nodeId)
					return lxerrors.New("flushing resource "+resourceId+" from node "+nodeId, err)
				}
			}
		}
	}
	return nil
}
