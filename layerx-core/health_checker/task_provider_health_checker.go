package health_checker

import (
	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/layerx/layerx-core/rpi_messenger"
	"github.com/emc-advanced-dev/layerx/layerx-core/tpi_messenger"
	"github.com/layer-x/layerx-commons/lxerrors"
	"time"
)

type HealthChecker struct {
	state *lxstate.State
}

func NewHealthChecker(state *lxstate.State) *HealthChecker {
	return &HealthChecker{
		state: state,
	}
}

func (hc *HealthChecker) FailDisconnectedTaskProviders() error {
	taskProviders, err := hc.state.TaskProviderPool.GetTaskProviders()
	if err != nil {
		return lxerrors.New("getting task providers from pool", err)
	}
	for _, taskProvider := range taskProviders {
		healthy, err := tpi_messenger.HealthCheck(hc.state.GetTpiUrl(), taskProvider)
		if err != nil {
			return lxerrors.New("performing health check on task provider "+taskProvider.Id, err)
		}
		if !healthy {
			logrus.WithFields(logrus.Fields{"task-provider": taskProvider}).Warnf("task provider disconnected")
			err = hc.state.TaskProviderPool.DeleteTaskProvider(taskProvider.Id)
			if err != nil {
				return lxerrors.New("removing failed task provider from active task provider pool", err)
			}
			taskProvider.TimeFailed = float64(time.Now().Unix())
			err = hc.state.FailedTaskProviderPool.AddTaskProvider(taskProvider)
			if err != nil {
				return lxerrors.New("adding failed task provider to failed task provider pool", err)
			}
			if taskProvider.FailoverTimeout == 0 {
				err = hc.destroyFailedTaskProvider(taskProvider)
				if err != nil {
					return lxerrors.New("destroying non-failover task provider that disconnected", err)
				}
			}
		}
	}
	return nil
}

func (hc *HealthChecker) ExpireTimedOutTaskProviders() error {
	failedTaskProviders, err := hc.state.FailedTaskProviderPool.GetTaskProviders()
	if err != nil {
		return lxerrors.New("getting task providers from failed pool", err)
	}
	for _, failedTaskProvider := range failedTaskProviders {
		expirationTime := failedTaskProvider.TimeFailed + failedTaskProvider.FailoverTimeout
		if float64(time.Now().Unix()) > expirationTime {
			logrus.WithFields(logrus.Fields{"task-provider": failedTaskProvider}).Warnf("failed-over task provider has expired, proceeding to purge")
			err = hc.destroyFailedTaskProvider(failedTaskProvider)
			if err != nil {
				return lxerrors.New("destroying non-failover task provider that disconnected", err)
			}
		}
	}
	return nil
}

func (hc *HealthChecker) destroyFailedTaskProvider(taskProvider *lxtypes.TaskProvider) error {
	allTasks, err := hc.state.GetAllTasks()
	if err != nil {
		return lxerrors.New("retrieving all tasks from state", err)
	}

	for taskId, task := range allTasks {
		if task.TaskProvider.Id == taskProvider.Id {
			logrus.WithFields(logrus.Fields{"task-provider": taskProvider, "task": task}).Debugf("destroying task for terminated task provider")
			for _, rpiUrl := range hc.state.GetRpiUrls() {
				err = rpi_messenger.SendKillTaskRequest(rpiUrl, taskId)
				if err != nil {
					logrus.WithFields(logrus.Fields{"err": err}).Warnf("sending kill task request to resource provider")
				}
			}
			taskPool, err := hc.state.GetTaskPoolContainingTask(taskId)
			if err != nil {
				return lxerrors.New("retrieving task pool containing task "+taskId, err)
			}
			err = taskPool.DeleteTask(taskId)
			if err != nil {
				return lxerrors.New("deleting task from task pool "+taskPool.GetKey(), err)
			}
		}
	}

	err = hc.state.FailedTaskProviderPool.DeleteTaskProvider(taskProvider.Id)
	if err != nil {
		return lxerrors.New("removing failed task provider from failed task provider pool", err)
	}

	return nil
}
