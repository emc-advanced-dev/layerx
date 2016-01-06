package task_launcher
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-core_v2/rpi_messenger"
	"fmt"
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
)

type TaskLauncher struct {
	rpiUrl string
	state *lxstate.State
}

func NewTaskLauncher(rpiUrl string, state *lxstate.State) *TaskLauncher {
	return &TaskLauncher{
		rpiUrl: rpiUrl,
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
		nodeId := task.SlaveId
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
		resourcesToUse := []*lxtypes.Resource{}
		for _, resource := range resourcesToUseMap {
			resourcesToUse = append(resourcesToUse, resource)
		}

		lxlog.Debugf(logrus.Fields{
			"tasks": fmt.Sprintf("%v",tasksToLaunch),
			"resources": fmt.Sprintf("%v", resourcesToUseMap),
			"node_id": fmt.Sprintf("%s",nodeId),
			"rpi_url": fmt.Sprintf("%s",tl.rpiUrl),
		}, "attempting to launch tasks on rpi")

		err = rpi_messenger.SendLaunchTasksMessage(tl.rpiUrl, tasksToLaunch, resourcesToUse)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"tasks": fmt.Sprintf("%v",tasksToLaunch),
				"resources": fmt.Sprintf("%v", resourcesToUseMap),
				"node_id": fmt.Sprintf("%s",nodeId),
				"rpi_url": fmt.Sprintf("%s",tl.rpiUrl),
			}, "trying to launch tasks on rpi")
			return lxerrors.New("sending launch task message to rpi", err)
		}
		//flush resources from node
		for resourceId, resource := range resourcesToUseMap {
			err := nodeResourcePool.DeleteResource(resourceId)
			if err != nil {
				lxlog.Errorf(logrus.Fields{
					"resource": fmt.Sprintf("%v",resource),
					"node_id": fmt.Sprintf("%s",nodeId),
				}, "flushing resource "+resourceId+" from node "+nodeId)
				return lxerrors.New("flushing resource "+resourceId+" from node "+nodeId, err)
			}
		}
	}
	return nil
}