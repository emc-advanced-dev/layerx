package health_checker

import (
	"github.com/layer-x/layerx-core_v2/lxstate"
"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/tpi_messenger"
	"github.com/layer-x/layerx-core_v2/lxtypes"
"github.com/layer-x/layerx-commons/lxlog"
"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-core_v2/rpi_messenger"
	"time"
)

type HealthChecker struct {
	tpiUrl string
	rpiUrl string
	state  *lxstate.State
}

func NewHealthChecker(tpiUrl, rpiUrl string, state *lxstate.State) *HealthChecker {
	return &HealthChecker{
		tpiUrl: tpiUrl,
		rpiUrl: rpiUrl,
		state: state,
	}
}

func (hc *HealthChecker) FailDisconnectedTaskProviders() error {
	taskProviders, err := hc.state.TaskProviderPool.GetTaskProviders()
	if err != nil {
		return lxerrors.New("getting task providers from pool", err)
	}
	for _, taskProvider := range taskProviders {
		healthy, err := tpi_messenger.HealthCheck(hc.tpiUrl, taskProvider)
		if err != nil {
			return lxerrors.New("performing health check on task provider "+taskProvider.Id, err)
		}
		if !healthy {
			lxlog.Warnf(logrus.Fields{"task-provider": taskProvider}, "task provider disconnected")
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
		lxlog.Warnf(logrus.Fields{"now": float64(time.Now().Unix()), "expire-time": expirationTime}, "did we expire yet?")
		if float64(time.Now().Unix()) > expirationTime {
			lxlog.Warnf(logrus.Fields{"task-provider": failedTaskProvider}, "failed-over task provider has expired, proceeding to purge")
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
			lxlog.Debugf(logrus.Fields{"task-provider": taskProvider, "task": task}, "destroying task for terminated task provider")
			err = rpi_messenger.SendKillTaskRequest(hc.rpiUrl, taskId)
			if err != nil {
				return lxerrors.New("sending kill task request to resource provider", err)
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