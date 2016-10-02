package framework_api_handlers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/mesos/mesos-go/mesosproto"
	"strings"
)

func HandleStatusUpdate(lxRpi *layerx_rpi_client.LayerXRpi, status *mesosproto.TaskStatus) error {
	taskState := status.GetState().String()
	statusString := "task " + status.GetTaskId().GetValue() + " in state" + taskState
	//ignore duplicate id bug
	//TODO: don't ignore this bug, but figure out where it's coming from
	if strings.Contains(taskState, "duplicate") {
		return nil
	}
	err := lxRpi.SubmitStatusUpdate(status)
	if err != nil {
		return errors.New("failed to submit status {"+statusString+"} to layerx core", err)
	}
	return nil
}
