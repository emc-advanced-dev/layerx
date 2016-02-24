package main_loop
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxerrors"
	"time"
	"github.com/layer-x/layerx-core_v2/tpi_messenger"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-core_v2/rpi_messenger"
	"github.com/layer-x/layerx-core_v2/task_launcher"
	"sync"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
)

var mainLoopLock = &sync.Mutex{}

//run as goroutine
func MainLoop(taskLauncher *task_launcher.TaskLauncher, state *lxstate.State, tpiUrl, rpiUrl string, driverErrc chan error) {
	for {
		errc := make(chan error)
		go func () {
			result := singleExeuction(state, taskLauncher, tpiUrl, rpiUrl)
			errc <- result
		}()
		err := <- errc
		if err != nil {
			driverErrc <- lxerrors.New("main loop failed while running", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func singleExeuction(state *lxstate.State, taskLauncher *task_launcher.TaskLauncher, tpiUrl, rpiUrl string) error {
	mainLoopLock.Lock()
	defer mainLoopLock.Unlock()
	taskProviderMap, err := state.TaskProviderPool.GetTaskProviders()
	if err != nil {
		return lxerrors.New("retrieving list of task providers from state", err)
	}
	taskProviders := []*lxtypes.TaskProvider{}
	for _, taskProvider := range taskProviderMap {
		taskProviders = append(taskProviders, taskProvider)
	}
	tpiErr := tpi_messenger.SendTaskCollectionMessage(tpiUrl, taskProviders)
	if tpiErr != nil {
		lxlog.Warnf(logrus.Fields{"error": err}, "failed sending task collection message to tpi. Is Tpi connected?")
	}

	rpiErr := rpi_messenger.SendResourceCollectionRequest(rpiUrl)
	if rpiErr != nil {
		lxlog.Warnf(logrus.Fields{"error": err}, "failed sending resource collection request to rpi. Is Rpi connected?")
	}

	if tpiErr != nil || rpiErr != nil {
		return nil
	}

	err = taskLauncher.LaunchStagedTasks()
	if err != nil {
		return lxerrors.New("launching staged tasks", err)
	}

	return nil
}