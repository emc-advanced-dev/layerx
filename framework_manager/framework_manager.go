package framework_manager
import (
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"strings"
	"net/http"
	"github.com/gogo/protobuf/proto"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"github.com/mesos/mesos-go/mesosproto"
	"fmt"
)

type frameworkManager struct {
	masterUpid *mesos_data.UPID
}

func NewFrameworkManager(masterUpid *mesos_data.UPID) *frameworkManager{
	return &frameworkManager{
		masterUpid: masterUpid,
	}
}

//notify a framework that it has successfully registered with the tpi
func (manager *frameworkManager) NotifyFrameworkRegistered(frameworkName, frameworkId, frameworkUpidString string) error {
	if frameworkName == "" {
		return lxerrors.New("framework must be named", nil)
	}
	frameworkUpid, err := mesos_data.UPIDFromString(frameworkUpidString)
	if err != nil {
		return lxerrors.New("converting upid string to upid struct", err)
	}

	masterState := &mesos_data.MesosState{
		Version: mesos_data.MESOS_VERSION,
		Leader: manager.masterUpid.String(),
	}
	masterInfo, err := masterState.ToMasterInfo()
	if err != nil {
		return lxerrors.New("converting master state to master info", err)
	}

	frameworkRegisteredMsg := &mesosproto.FrameworkRegisteredMessage{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String(frameworkId),
		},
		MasterInfo: masterInfo,
	}
	resp, _, err := manager.sendMessage(frameworkUpid, frameworkRegisteredMsg, "/mesos.internal.FrameworkRegisteredMessage")
	if err != nil {
		return lxerrors.New("sending registered message to framework", err)
	}
	if !(resp.StatusCode == 200 || resp.StatusCode == 202) {
		statusCode := fmt.Sprintf("%v", resp.StatusCode)
		return lxerrors.New("expected 200 or 202 response from framework, got "+statusCode, nil)
	}
	return nil
}


func (manager *frameworkManager) sendMessage(destination *mesos_data.UPID, message proto.Message, path string) (*http.Response, []byte, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url := destination.Host +":"+ destination.Port
	path = "/" + destination.ID + path
	headers := map[string]string{
		"Libprocess-From": manager.masterUpid.String(),
		"Content-Type":    "application/json",
	}
	resp, data, err := lxhttpclient.Post(url, path, headers, message)
	if err != nil {
		err = lxerrors.New("sending data("+string(data)+") to framework", err)
	}
	return resp, data, err
}