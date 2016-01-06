package main_loop
import (
	"github.com/layer-x/layerx-core_v2/lxstate"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/layer-x/layerx-commons/lxerrors"
	"time"
	"github.com/layer-x/layerx-core_v2/tpi_messenger"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-core_v2/rpi_messenger"
	"github.com/layer-x/layerx-core_v2/task_launcher"
)

//run as goroutine
func MainLoop(actionQueue lxactionqueue.ActionQueue, taskLauncher *task_launcher.TaskLauncher, state *lxstate.State, tpiUrl, rpiUrl string, driverErrc chan error) {
	for {
		errc := make(chan error)
		actionQueue.Push(func () {
			result := singleExeuction(state, taskLauncher, tpiUrl, rpiUrl)
			errc <- result
		})
		err := <- errc
		if err != nil {
			driverErrc <- lxerrors.New("main loop failed while running", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func singleExeuction(state *lxstate.State, taskLauncher *task_launcher.TaskLauncher, tpiUrl, rpiUrl string) error {
	taskProviderMap, err := state.TaskProviderPool.GetTaskProviders()
	if err != nil {
		return lxerrors.New("retrieving list of task providers from state", err)
	}
	taskProviders := []*lxtypes.TaskProvider{}
	for _, taskProvider := range taskProviderMap {
		taskProviders = append(taskProviders, taskProvider)
	}
	err = tpi_messenger.SendTaskCollectionMessage(tpiUrl, taskProviders)
	if err != nil {
		return lxerrors.New("sending task collection message to tpi", err)
	}

	err = rpi_messenger.SendResourceCollectionRequest(rpiUrl)
	if err != nil {
		return lxerrors.New("sending resource collection request to rpi", err)
	}
	err = taskLauncher.LaunchStagedTasks()
	if err != nil {
		return lxerrors.New("launching staged tasks", err)
	}

	return nil
}