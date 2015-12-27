package mesos_api_helpers
import (
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-commons/lxerrors"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
)

func HandleLaunchTasksRequest(tpi *layerx_tpi.LayerXTpi, frameworkId string, mesosTasks []*mesosproto.TaskInfo) error {
	for _, mesosTask := range mesosTasks {
		lxTask := lxtypes.NewTaskFromMesos(mesosTask)
		err := tpi.SubmitTask(frameworkId, lxTask)
		if err != nil {
			lxlog.Errorf(logrus.Fields{
				"error": err.Error(),
				"frameworkId": frameworkId,
				"tpi": tpi,
				"task": lxTask,
			}, "submitting task to layer-x core")
			return lxerrors.New("submitting task "+lxTask.TaskId+" to layer-x core", err)
		}
	}
	return nil
}