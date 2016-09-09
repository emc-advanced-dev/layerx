package layerx_rpi_client
import "github.com/emc-advanced-dev/layerx-core/lxtypes"

type RpiInfo struct {
	Name string `json:"name"`
	Url string `json:"rpi_url"`
}

type LaunchTasksMessage struct {
	TasksToLaunch []*lxtypes.Task `json:"tasks_to_launch"`
	ResourcesToUse []*lxtypes.Resource `json:"resources_to_use"`
}
