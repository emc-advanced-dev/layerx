package layerx_rpi_client
import (
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/layer-x/layerx-commons/lxerrors"
	"fmt"
	"github.com/mesos/mesos-go/mesosproto"
)


type LayerXRpi struct {
	CoreURL string
}

const (
	SubmitResource             = "/SubmitResource"
	SubmitStatusUpdate         = "/SubmitStatusUpdate"
)

//call this method when submitting
// a new resource from the rpi
func (rpi *LayerXRpi) SubmitResource(resource *lxtypes.Resource) error {
	resp, _, err := lxhttpclient.Post(rpi.CoreURL, SubmitResource, nil, resource)
	if err != nil {
		return lxerrors.New("POSTing resource to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing resource to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}

//call this method when submitting
// a status update from the rpi
func (rpi *LayerXRpi) SubmitStatusUpdate(status *mesosproto.TaskStatus) error {
	resp, _, err := lxhttpclient.Post(rpi.CoreURL, SubmitStatusUpdate, nil, status)
	if err != nil {
		return lxerrors.New("POSTing TaskStatus to LayerX core server", err)
	}
	if resp.StatusCode != 202 {
		msg := fmt.Sprintf("POSTing TaskStatus to LayerX core server; status code was %v, expected 202", resp.StatusCode)
		return lxerrors.New(msg, err)
	}
	return nil
}
