package mesos_master_api
import (
"fmt"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-commons/lxmartini"
"net/http"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/handlers"
	"github.com/layer-x/layerx-commons/lxerrors"
)

const (
	GET_MASTER_STATE                      = "/master/state.json"
	GET_MASTER_STATE_DEPRECATED           = "/state.json"
	MESOS_SCHEDULER_CALL                  = "/master/mesos.scheduler.Call"
	REGISTER_FRAMEWORK_MESSAGE            = "/master/mesos.internal.RegisterFrameworkMessage"
	REREGISTER_FRAMEWORK_MESSAGE          = "/master/mesos.internal.ReregisterFrameworkMessage"
	UNREGISTER_FRAMEWORK_MESSAGE          = "/master/mesos.internal.UnregisterFrameworkMessage"
	LAUNCH_TASKS_MESSAGE                  = "/master/mesos.internal.LaunchTasksMessage"
	RECONCILE_TASKS_MESSAGE               = "/master/mesos.internal.ReconcileTasksMessage"
	KILL_TASK_MESSAGE                     = "/master/mesos.internal.KillTaskMessage"
	STATUS_UPDATE_ACKNOWLEDGEMENT_MESSAGE = "/master/mesos.internal.StatusUpdateAcknowledgementMessage"
	REVIVE_OFFERS_MESSAGE                 = "/master/mesos.internal.ReviveOffersMessage"
)

func RunMasterServer(port int, masterUpid string, errc chan error) {
	portStr := fmt.Sprintf(":%v", port)
	lxlog.Infof(logrus.Fields{
		"port": port,
	}, "Master Server initialized")

	m := lxmartini.QuietMartini()

	m.Get(GET_MASTER_STATE, func(res http.ResponseWriter) {
		data, err := handlers.GetMesosState(masterUpid)
		if err != nil {
			res.WriteHeader(500)
			lxlog.Errorf(logrus.Fields{
				"port": port,
			}, "Retreiving master state: %s", err.Error())
			errc <- lxerrors.New("retreiving master state", err)
			return
		}
		res.Write(data)
	})

	m.Get(GET_MASTER_STATE_DEPRECATED, func(res http.ResponseWriter) {
		data, err := handlers.GetMesosState(masterUpid)
		if err != nil {
			res.WriteHeader(500)
			lxlog.Errorf(logrus.Fields{
				"port": port,
			}, "Retreiving master state: %s", err.Error())
			errc <- lxerrors.New("retreiving master state", err)
			return
		}
		res.Write(data)
	})



	m.RunOnAddr(portStr)
}