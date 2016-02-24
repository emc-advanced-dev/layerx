package mesos_framework_api
import (
	"github.com/mesos/mesos-go/scheduler"
	"github.com/mesos/mesos-go/mesosproto"
"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-mesos-rpi_v2/mesos_framework_api/framework_api_handlers"
)

type MesosScheduler interface {
	scheduler.Scheduler
	GetDriver() scheduler.SchedulerDriver
}

type rpiMesosScheduler struct {
	driver scheduler.SchedulerDriver
	driverc chan scheduler.SchedulerDriver
	lxRpi *layerx_rpi_client.LayerXRpi
}

func NewRpiMesosScheduler(lxRpi *layerx_rpi_client.LayerXRpi) *rpiMesosScheduler {
	return &rpiMesosScheduler{
		driver: nil,
		driverc: make(chan scheduler.SchedulerDriver),
		lxRpi: lxRpi,
	}
}

func (s *rpiMesosScheduler) GetDriver() scheduler.SchedulerDriver {
	if s.driver == nil {
		s.driver = <- s.driverc
	}
	return s.driver
}

func (s *rpiMesosScheduler) Registered(driver scheduler.SchedulerDriver, frameworkId *mesosproto.FrameworkID, masterInfo *mesosproto.MasterInfo) {
	lxlog.Infof(logrus.Fields{
		"framework_id": frameworkId.GetValue(),
		"master_id":    masterInfo.GetId(),
	}, "Scheduler Registered to Master %v\n", masterInfo)
	s.driverc <- driver
}

func (s *rpiMesosScheduler) Reregistered(driver scheduler.SchedulerDriver, masterInfo *mesosproto.MasterInfo) {
	lxlog.Infof(logrus.Fields{}, "Scheduler Re-Registered with Master %v\n", masterInfo)
	s.driverc <- driver
}

func (s *rpiMesosScheduler) Disconnected(scheduler.SchedulerDriver) {
	lxlog.Infof(logrus.Fields{}, "Scheduler Disconnected")
}

func (s *rpiMesosScheduler) ResourceOffers(driver scheduler.SchedulerDriver, offers []*mesosproto.Offer) {
	lxlog.Infof(logrus.Fields{}, "Collecting %v offers from Mesos Master...\n", len(offers))
	go func(){
		err := framework_api_handlers.HandleResourceOffers(s.lxRpi, offers)
		if err != nil {
			lxlog.Fatalf(logrus.Fields{
				"error": err,
			}, "handling resource offers from mesos master")
		}
	}()
}

func (s *rpiMesosScheduler) StatusUpdate(driver scheduler.SchedulerDriver, status *mesosproto.TaskStatus) {
	lxlog.Infof(logrus.Fields{}, "Status update: task "+status.GetTaskId().GetValue()+" is in state "+status.State.Enum().String()+" with message %s", status.GetMessage())
	go func(){
		err := framework_api_handlers.HandleStatusUpdate(s.lxRpi, status)
		if err != nil {
			lxlog.Fatalf(logrus.Fields{
				"error": err,
			}, "handling task status update from mesos master")
		}
	}()
}

func (s *rpiMesosScheduler) OfferRescinded(driver scheduler.SchedulerDriver, id *mesosproto.OfferID) {
	lxlog.Infof(logrus.Fields{}, "Offer '%v' rescinded.\n", *id)
}

func (s *rpiMesosScheduler) FrameworkMessage(driver scheduler.SchedulerDriver, exId *mesosproto.ExecutorID, slvId *mesosproto.SlaveID, msg string) {
	lxlog.Infof(logrus.Fields{}, "Received framework message from executor '%v' on slave '%v': %s.\n", *exId, *slvId, msg)
}

func (s *rpiMesosScheduler) SlaveLost(driver scheduler.SchedulerDriver, id *mesosproto.SlaveID) {
	lxlog.Infof(logrus.Fields{}, "Slave '%v' lost.\n", *id)
}

func (s *rpiMesosScheduler) ExecutorLost(driver scheduler.SchedulerDriver, exId *mesosproto.ExecutorID, slvId *mesosproto.SlaveID, i int) {
	lxlog.Infof(logrus.Fields{}, "Executor '%v' lost on slave '%v' with exit code: %v.\n", *exId, *slvId, i)
}

func (s *rpiMesosScheduler) Error(driver scheduler.SchedulerDriver, err string) {
	lxlog.Infof(logrus.Fields{}, "Scheduler received error:", err)
}
