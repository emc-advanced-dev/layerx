package layerx_rpi_client
import "github.com/layer-x/layerx-core_v2/lxtypes"

type LaunchTasksMessage struct {
	TasksToLaunch []*lxtypes.Task `json:"tasks_to_launch"`
	ResourcesToUse []*lxtypes.Resource `json:"resources_to_use"`
}