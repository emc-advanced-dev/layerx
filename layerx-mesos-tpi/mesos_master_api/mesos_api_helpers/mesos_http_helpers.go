package mesos_api_helpers

import (
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/mesos_master_api/mesos_data"
	"github.com/emc-advanced-dev/pkg/errors"
)

func ProcessMesosHttpRequest(req *http.Request) (*mesos_data.UPID, []byte, int, error) {
	data, err := ioutil.ReadAll(req.Body)
	if req.Body != nil {
		defer req.Body.Close()
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Errorf("could not read  REGISTER_FRAMEWORK_MESSAGE request body")
		return nil, empty, 500, errors.New("could not read  mesos http request body", err)
	}
	requestingFramework := req.Header.Get("Libprocess-From")
	if requestingFramework == "" {
		logrus.Errorf("missing required header: %s", "Libprocess-From")
		return nil, empty, 400, nil
	}
	upid, err := mesos_data.UPIDFromString(requestingFramework)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Errorf("could not parse pid of requesting framework")
		return nil, empty, 500, errors.New("could not parse pid of requesting framework", err)
	}
	return upid, data, -1, nil
}
