package mesos_framework_api

import (
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-mesos-rpi_v2/mesos_framework_api/framework_api_handlers"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
)

type MesosScheduler interface {
	scheduler.Scheduler
	GetDriver() scheduler.SchedulerDriver
}

type rpiMesosScheduler struct {
	driver  scheduler.SchedulerDriver
	driverc chan scheduler.SchedulerDriver
	lxRpi   *layerx_rpi_client.LayerXRpi
}

func NewRpiMesosScheduler(lxRpi *layerx_rpi_client.LayerXRpi) *rpiMesosScheduler {
	return &rpiMesosScheduler{
		driver:  nil,
		driverc: make(chan scheduler.SchedulerDriver),
		lxRpi:   lxRpi,
	}
}

func (s *rpiMesosScheduler) GetDriver() scheduler.SchedulerDriver {
	if s.driver == nil {
		s.driver = <-s.driverc
	}
	return s.driver
}

func (s *rpiMesosScheduler) Registered(driver scheduler.SchedulerDriver, frameworkId *mesosproto.FrameworkID, masterInfo *mesosproto.MasterInfo) {
	logrus.WithFields(logrus.Fields{
		"framework_id": frameworkId.GetValue(),
		"master_id":    masterInfo.GetId(),
	}).Infof("Scheduler Registered to Master %v\n", masterInfo)
	s.driverc <- driver
}

func (s *rpiMesosScheduler) Reregistered(driver scheduler.SchedulerDriver, masterInfo *mesosproto.MasterInfo) {
	logrus.Infof("Scheduler Re-Registered with Master %v\n", masterInfo)
	s.driverc <- driver
}

func (s *rpiMesosScheduler) Disconnected(scheduler.SchedulerDriver) {
	logrus.Infof("Scheduler Disconnected")
}

func (s *rpiMesosScheduler) ResourceOffers(driver scheduler.SchedulerDriver, offers []*mesosproto.Offer) {
	logrus.Infof("Collecting %v offers from Mesos Master...\n", len(offers))
	go func() {
		err := framework_api_handlers.HandleResourceOffers(s.lxRpi, offers)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("handling resource offers from mesos master")
		}
	}()
}

func (s *rpiMesosScheduler) StatusUpdate(driver scheduler.SchedulerDriver, status *mesosproto.TaskStatus) {
	logrus.Infof("Status update: task "+status.GetTaskId().GetValue()+" is in state "+status.State.Enum().String()+" with message %s", status.GetMessage())
	go func() {
		err := framework_api_handlers.HandleStatusUpdate(s.lxRpi, status)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("handling task status update from mesos master")
		}
	}()
}

func (s *rpiMesosScheduler) OfferRescinded(driver scheduler.SchedulerDriver, id *mesosproto.OfferID) {
	logrus.Infof("Offer '%v' rescinded.\n", *id)
}

func (s *rpiMesosScheduler) FrameworkMessage(driver scheduler.SchedulerDriver, exId *mesosproto.ExecutorID, slvId *mesosproto.SlaveID, msg string) {
	logrus.Infof("Received framework message from executor '%v' on slave '%v': %s.\n", *exId, *slvId, msg)
}

func (s *rpiMesosScheduler) SlaveLost(driver scheduler.SchedulerDriver, id *mesosproto.SlaveID) {
	logrus.Infof("Slave '%v' lost.\n", *id)
}

func (s *rpiMesosScheduler) ExecutorLost(driver scheduler.SchedulerDriver, exId *mesosproto.ExecutorID, slvId *mesosproto.SlaveID, i int) {
	logrus.Infof("Executor '%v' lost on slave '%v' with exit code: %v.\n", *exId, *slvId, i)
}

func (s *rpiMesosScheduler) Error(driver scheduler.SchedulerDriver, err string) {
	logrus.Infof("Scheduler received error:", err)
}
