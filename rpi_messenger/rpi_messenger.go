package rpi_messenger
import (
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"fmt"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
)

const (
	COLLECT_RESOURCES = "/collect_resources"
	LAUNCH_TASKS = "/launch_resources"
	KILL_TASK = "/kill_task"
)

func SendResourceCollectionRequest(tpiUrl string) error {
	resp, _, err := lxhttpclient.Post(tpiUrl, COLLECT_RESOURCES, nil, nil)
	if err != nil {
		return lxerrors.New("POSTing COLLECT_RESOURCES to TPI server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing COLLECT_RESOURCES to TPI server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}

func SendLaunchTasksMessage(tpiUrl string, tasksToLaunch []*lxtypes.Task, resourcesToUse []*lxtypes.Resource) error {
	launchTasksMessage := &layerx_rpi_client.LaunchTasksMessage{
		TasksToLaunch: tasksToLaunch,
		ResourcesToUse: resourcesToUse,
	}
	resp, _, err := lxhttpclient.Post(tpiUrl, LAUNCH_TASKS, nil, launchTasksMessage)
	if err != nil {
		return lxerrors.New("POSTing tasksToLaunch to TPI server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing tasksToLaunch to TPI server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}

func SendKillTaskRequest(tpiUrl string, taskId string) error {
	resp, _, err := lxhttpclient.Post(tpiUrl, KILL_TASK+"/"+taskId, nil, nil)
	if err != nil {
		return lxerrors.New("POSTing KillTask request for task "+taskId+" to TPI server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing KillTask request for task "+taskId+" to TPI server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}