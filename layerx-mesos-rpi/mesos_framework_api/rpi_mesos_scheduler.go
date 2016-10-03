package mesos_framework_api

import (
	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-rpi/mesos_framework_api/framework_api_handlers"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
)

type MesosScheduler interface {
	scheduler.Scheduler
	GetDriver() scheduler.SchedulerDriver
}

type rpiMesosScheduler struct {
	driver        scheduler.SchedulerDriver
	driverc       chan scheduler.SchedulerDriver
	core          *layerx_rpi_client.LayerXRpi
	TaskChan      chan *lxtypes.Task
	taskQueue     []*lxtypes.Task
	tasksLaunched int
}

func NewRpiMesosScheduler(lxRpi *layerx_rpi_client.LayerXRpi) *rpiMesosScheduler {
	s := &rpiMesosScheduler{
		driver:  nil,
		driverc: make(chan scheduler.SchedulerDriver),
		core:    lxRpi,
		TaskChan: make(chan *lxtypes.Task),
	}
	//process tasks from the chan into the queue
	go func(){
		for {
			task := <-s.TaskChan
			logrus.Debugf("popping task %v", task)
			s.taskQueue = append(s.taskQueue, task)
		}
	}()
	return s
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
	logrus.Debugf("Searching %v offers from Mesos Master to place %v tasks... %v launched so far", len(offers), len(s.taskQueue), s.tasksLaunched)
	if s.tasksLaunched < len(s.taskQueue) {
		task := s.taskQueue[s.tasksLaunched]
		offersToUse := []*mesosproto.OfferID{}
		for _, offer := range offers {
			if offer.GetSlaveId().GetValue() == task.NodeId {
				offersToUse = append(offersToUse, offer.GetId())
			} else {
				logrus.Debugf("declining offer %v", offer.GetId().GetValue())
				driver.DeclineOffer(offer.GetId(), &mesosproto.Filters{})
				driver.ReviveOffers()
			}
		}
		if status, err := driver.LaunchTasks(offersToUse, []*mesosproto.TaskInfo{task.ToMesos()}, &mesosproto.Filters{}); err != nil || status != mesosproto.Status_DRIVER_RUNNING {
			logrus.Errorf("failed to launch task %v on offers %v\n failed with status %v and error %v", task.ToMesos(), offersToUse, status, err)
			return
		}
		logrus.Infof("Successfully launched %v", task)
		s.tasksLaunched++
	} else {
		//nothing to do, decline all offers
		for _, offer := range offers {
			logrus.Debugf("declining offer %v", offer.GetId().GetValue())
			driver.DeclineOffer(offer.GetId(), &mesosproto.Filters{})
			driver.ReviveOffers()
		}
	}
}

func (s *rpiMesosScheduler) StatusUpdate(driver scheduler.SchedulerDriver, status *mesosproto.TaskStatus) {
	logrus.Infof("Status update: task " + status.GetTaskId().GetValue() + " is in state " + status.State.Enum().String() + " with message %s", status.GetMessage())
	go func() {
		err := framework_api_handlers.HandleStatusUpdate(s.core, status)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("handling task status update from mesos master")
		}
	}()
}

func (s *rpiMesosScheduler) OfferRescinded(driver scheduler.SchedulerDriver, id *mesosproto.OfferID) {
	logrus.Infof("Offer '%v' rescinded. Notifying Core and declining / reviving offer.\n", *id)
	if err := framework_api_handlers.HandleOfferRescinded(s.core, driver, id); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Errorf("handling offer rescinded message from mesos master")
	}

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
