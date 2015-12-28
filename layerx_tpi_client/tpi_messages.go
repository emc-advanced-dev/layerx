package layerx_tpi
import (
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/mesos/mesos-go/mesosproto"
)

type TpiRegistrationMessage struct {
	TpiUrl string `json:"tpi_url"`
}

type CollectTasksMessage struct {
	TaskProivders []*lxtypes.TaskProvider `json:"task_providers"`
}

type UpdateTaskStatusMessage struct {
	TaskProvider *lxtypes.TaskProvider `json:"task_provider"`
	TaskStatus *mesosproto.TaskStatus `json:"task_status"`
}