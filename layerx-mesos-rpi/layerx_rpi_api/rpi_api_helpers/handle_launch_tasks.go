package rpi_api_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/Sirupsen/logrus"
)

func LaunchTasks(taskQueue chan *lxtypes.Task, launchTasksMessage layerx_rpi_client.LaunchTasksMessage) error {
	resources := launchTasksMessage.ResourcesToUse
	resourceCount := len(resources)
	if resourceCount < 1 {
		return errors.New("need at least one resource to launch a task", nil)
	}
	var index int
	tasks := launchTasksMessage.TasksToLaunch
	for _, task := range tasks {
		logrus.Debug("pushing task ", task, " onto queue")
		//select any node in the resource list to use
		//make sure node is set to the target slave
		task.NodeId = resources[index%resourceCount].NodeId
		taskQueue <- task
		logrus.Debug("queue is len", len(taskQueue))
		index++
	}
	return nil
}
