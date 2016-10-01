package rpi_messenger

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
)

const (
	COLLECT_RESOURCES = "/collect_resources"
	LAUNCH_TASKS      = "/launch_tasks"
	KILL_TASK         = "/kill_task"
)

func SendResourceCollectionRequest(rpiUrl string) error {
	logrus.Debug("sending resource collection message to", rpiUrl)
	resp, _, err := lxhttpclient.Post(rpiUrl, COLLECT_RESOURCES, nil, nil)
	if err != nil {
		return errors.New("POSTing COLLECT_RESOURCES to RPI server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing COLLECT_RESOURCES to RPI server; status code was %v, expected 202", resp.StatusCode)
		return errors.New(msg, err)
	}
	return nil
}

func SendLaunchTasksMessage(rpiUrl string, tasksToLaunch []*lxtypes.Task, resourcesToUse []*lxtypes.Resource) error {
	launchTasksMessage := &layerx_rpi_client.LaunchTasksMessage{
		TasksToLaunch:  tasksToLaunch,
		ResourcesToUse: resourcesToUse,
	}
	resp, _, err := lxhttpclient.Post(rpiUrl, LAUNCH_TASKS, nil, launchTasksMessage)
	if err != nil {
		return errors.New("POSTing tasksToLaunch to RPI server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing tasksToLaunch to RPI server; status code was %v, expected 202", resp.StatusCode)
		return errors.New(msg, err)
	}
	return nil
}

func SendKillTaskRequest(rpiUrl string, taskId string) error {
	resp, _, err := lxhttpclient.Post(rpiUrl, KILL_TASK+"/"+taskId, nil, nil)
	if err != nil {
		return errors.New("POSTing KillTask request for task "+taskId+" to RPI server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing KillTask request for task "+taskId+" to RPI server; status code was %v, expected 202", resp.StatusCode)
		return errors.New(msg, err)
	}
	return nil
}
