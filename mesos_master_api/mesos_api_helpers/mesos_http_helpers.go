package mesos_api_helpers
import (
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/Sirupsen/logrus"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
"io/ioutil"
"net/http"
)

func ProcessMesosHttpRequest(req *http.Request) (*mesos_data.UPID, []byte, int, error) {
	data, err := ioutil.ReadAll(req.Body)
	if req.Body != nil {
		defer req.Body.Close()
	}
	if err != nil {
		lxlog.Errorf(logrus.Fields{
			"error": err,
		}, "could not read  REGISTER_FRAMEWORK_MESSAGE request body")
		return nil, empty, 500, lxerrors.New("could not read  mesos http request body", err)
	}
	requestingFramework := req.Header.Get("Libprocess-From")
	if requestingFramework == "" {
		lxlog.Errorf(logrus.Fields{}, "missing required header: %s", "Libprocess-From")
		return nil, empty, 400, nil
	}
	upid, err := mesos_data.UPIDFromString(requestingFramework)
	if err != nil {
		lxlog.Errorf(logrus.Fields{
			"error": err,
		}, "could not parse pid of requesting framework")
		return nil, empty, 500, lxerrors.New("could not parse pid of requesting framework", err)
	}
	return upid, data, -1, nil
}