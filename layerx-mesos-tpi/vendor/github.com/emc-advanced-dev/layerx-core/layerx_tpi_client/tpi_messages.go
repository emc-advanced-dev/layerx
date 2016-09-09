package layerx_tpi_client

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/mesos/mesos-go/mesosproto"
)

type TpiRegistrationMessage struct {
	TpiUrl string `json:"tpi_url"`
}

type CollectTasksMessage struct {
	TaskProviders []*lxtypes.TaskProvider `json:"task_providers"`
}

type UpdateTaskStatusMessage struct {
	TaskProvider *lxtypes.TaskProvider  `json:"task_provider"`
	TaskStatus   *mesosproto.TaskStatus `json:"task_status"`
}

type HealthCheckTaskProviderMessage struct {
	TaskProvider *lxtypes.TaskProvider `json:"task_provider"`
}
