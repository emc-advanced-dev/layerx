package main_loop

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/health_checker"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/layerx/layerx-core/rpi_messenger"
	"github.com/emc-advanced-dev/layerx/layerx-core/task_launcher"
	"github.com/emc-advanced-dev/layerx/layerx-core/tpi_messenger"
	"github.com/emc-advanced-dev/pkg/errors"
)

var mainLoopLock = &sync.Mutex{}

//run as goroutine
func MainLoop(taskLauncher *task_launcher.TaskLauncher, healthChecker *health_checker.HealthChecker, state *lxstate.State, driverErrc chan error) {
	for {
		errc := make(chan error)
		go func() {
			result := singleExeuction(state, taskLauncher, healthChecker)
			errc <- result
		}()
		err := <-errc
		if err != nil {
			driverErrc <- errors.New("main loop failed while running", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func singleExeuction(state *lxstate.State, taskLauncher *task_launcher.TaskLauncher, healthChecker *health_checker.HealthChecker) error {
	mainLoopLock.Lock()
	defer mainLoopLock.Unlock()
	taskProviderMap, err := state.TaskProviderPool.GetTaskProviders()
	if err != nil {
		return errors.New("retrieving list of task providers from state", err)
	}
	taskProviders := []*lxtypes.TaskProvider{}
	for _, taskProvider := range taskProviderMap {
		taskProviders = append(taskProviders, taskProvider)
	}
	tpiErr := tpi_messenger.SendTaskCollectionMessage(state.GetTpiUrl(), taskProviders)
	if tpiErr != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Warnf("failed sending task collection message to tpi. Is Tpi connected?")
	}
	var rpiErr error
	for _, rpiUrl := range state.GetRpiUrls() {
		rpiErr = rpi_messenger.SendResourceCollectionRequest(rpiUrl)
		if rpiErr != nil {
			logrus.WithFields(logrus.Fields{"error": err}).Warnf("failed sending resource collection request to rpi. Is Rpi connected?")
		}
	}

	if tpiErr != nil || rpiErr != nil {
		return nil
	}

	err = healthChecker.FailDisconnectedTaskProviders()
	if err != nil {
		return errors.New("failing disconnected task providers", err)
	}

	err = healthChecker.ExpireTimedOutTaskProviders()
	if err != nil {
		return errors.New("expiring timed out providers", err)
	}

	err = taskLauncher.LaunchStagedTasks()
	if err != nil {
		return errors.New("launching staged tasks", err)
	}

	return nil
}
