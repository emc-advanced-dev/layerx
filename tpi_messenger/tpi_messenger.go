package tpi_messenger
import (
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"fmt"
)

const (
	COLLECT_TASKS = "/collect_tasks"
	UPDATE_TASK_STATUS = "/update_task_status"
)

func SendTaskCollectionMessage(tpiUrl string, taskProviders []*lxtypes.TaskProvider) error {
	collectTasksMessage := &layerx_tpi_client.CollectTasksMessage{
		TaskProviders: taskProviders,
	}
	resp, _, err := lxhttpclient.Post(tpiUrl, COLLECT_TASKS, nil, collectTasksMessage)
	if err != nil {
		return lxerrors.New("POSTing CollectTasksMessage to TPI server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing CollectTasksMessage to TPI server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}

func SendStatusUpdate(tpiUrl string, taskProvider *lxtypes.TaskProvider, status *mesosproto.TaskStatus) error {
	updateTaskStatusMessage := &layerx_tpi_client.UpdateTaskStatusMessage{
		TaskProvider: taskProvider,
		TaskStatus: status,
	}
	resp, _, err := lxhttpclient.Post(tpiUrl, UPDATE_TASK_STATUS, nil, updateTaskStatusMessage)
	if err != nil {
		return lxerrors.New("POSTing TaskStatus to TPI server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing TaskStatus to TPI server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}