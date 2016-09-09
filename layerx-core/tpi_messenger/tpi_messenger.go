package tpi_messenger

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/mesos/mesos-go/mesosproto"

	"net/http"
)

const (
	COLLECT_TASKS              = "/collect_tasks"
	UPDATE_TASK_STATUS         = "/update_task_status"
	HEALTH_CHECK_TASK_PROVIDER = "/health_check_task_provider"
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
		logrus.WithFields(logrus.Fields{"response": resp}).Warnf("error: " + msg)
	}
	return nil
}

func SendStatusUpdate(tpiUrl string, taskProvider *lxtypes.TaskProvider, status *mesosproto.TaskStatus) error {
	updateTaskStatusMessage := &layerx_tpi_client.UpdateTaskStatusMessage{
		TaskProvider: taskProvider,
		TaskStatus:   status,
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

func HealthCheck(tpiUrl string, taskProvider *lxtypes.TaskProvider) (bool, error) {
	healthCheckTaskProvider := &layerx_tpi_client.HealthCheckTaskProviderMessage{
		TaskProvider: taskProvider,
	}
	resp, _, err := lxhttpclient.Post(tpiUrl, HEALTH_CHECK_TASK_PROVIDER, nil, healthCheckTaskProvider)
	if err != nil {
		return false, lxerrors.New("POSTing HealthCheckTaskProviderMessage to TPI server", err)
	}
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusGone {
		return false, nil
	}
	msg := fmt.Sprintf("POSTing HealthCheckTaskProviderMessage to TPI server; status code was %v, expected 200 or 410", resp.StatusCode)
	return false, lxerrors.New(msg, err)
}
