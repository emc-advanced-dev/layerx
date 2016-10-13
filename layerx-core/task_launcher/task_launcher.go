package task_launcher

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/layerx/layerx-core/rpi_messenger"
	"github.com/emc-advanced-dev/pkg/errors"
	"sync"
	"time"
)

var retryInterval = time.Second * 15

type TaskLauncher struct {
	state            *lxstate.State
	recentlyLaunched map[string]bool
	mapLock          sync.RWMutex
}

func NewTaskLauncher(state *lxstate.State) *TaskLauncher {
	return &TaskLauncher{
		state: state,
		recentlyLaunched: make(map[string]bool),
	}
}

func (tl *TaskLauncher) LaunchStagedTasks() error {
	stagingTasks, err := tl.state.StagingTaskPool.GetTasks()
	if err != nil {
		return errors.New("retrieving staging tasks", err)
	}

	nodeTaskMap := make(map[string][]*lxtypes.Task)
	for _, task := range stagingTasks {
		//don't rapidly try to relaunch tasks
		tl.mapLock.RLock()
		if tl.recentlyLaunched[task.TaskId] {
			tl.mapLock.RUnlock()
			continue
		}
		tl.mapLock.RUnlock()
		nodeId := task.NodeId
		_, ok := nodeTaskMap[nodeId]
		if !ok {
			nodeTaskMap[nodeId] = []*lxtypes.Task{}
		}
		//mark task so we don't try to relaunch it for the retryInterval
		tl.mapLock.Lock()
		tl.recentlyLaunched[task.TaskId] = true
		tl.mapLock.Unlock()
		go func(){
			time.Sleep(retryInterval)
			tl.mapLock.Lock()
			delete(tl.recentlyLaunched, task.TaskId)
			tl.mapLock.Unlock()
		}()
		nodeTaskMap[nodeId] = append(nodeTaskMap[nodeId], task)
	}
	for nodeId, tasksToLaunch := range nodeTaskMap {
		nodeResourcePool, err := tl.state.NodePool.GetNodeResourcePool(nodeId)
		if err != nil {
			return errors.New("finding resource pool for node "+nodeId, err)
		}
		resourcesToUseMap, err := nodeResourcePool.GetResources()
		if err != nil {
			return errors.New("retrieving resource list for node "+nodeId, err)
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
				return errors.New("retreiving rpi for name", err)
			}
			err = rpi_messenger.SendLaunchTasksMessage(rpi.Url, tasksToLaunch, resourcesToUse)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"tasks":     fmt.Sprintf("%v", tasksToLaunch),
					"resources": fmt.Sprintf("%v", resourcesToUseMap),
					"node_id":   fmt.Sprintf("%s", nodeId),
					"rpi_url":   rpi.Url,
				}).Errorf("trying to launch tasks on rpi")
				return errors.New("sending launch task message to rpi", err)
			}
			//flush resources from node
			for resourceId, resource := range resourcesToUseMap {
				err := nodeResourcePool.DeleteResource(resourceId)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"resource": fmt.Sprintf("%v", resource),
						"node_id":  fmt.Sprintf("%s", nodeId),
					}).Errorf("flushing resource " + resourceId + " from node " + nodeId)
					return errors.New("flushing resource "+resourceId+" from node "+nodeId, err)
				}
			}
		}
	}
	return nil
}
