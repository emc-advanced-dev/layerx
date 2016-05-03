package mesos_master_api_test

import (
	. "github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"encoding/json"
	"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
	"github.com/layer-x/layerx-mesos-tpi_v2/fakes"
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	core_fakes "github.com/layer-x/layerx-core_v2/fakes"
	"github.com/layer-x/layerx-commons/lxmartini"
	"fmt"
	"github.com/mesos/mesos-go/mesosproto"
)

var _ = Describe("MasterApiServer", func() {
	fakeMasterUpid, _ := mesos_data.UPIDFromString("master@127.0.0.1:3031")
	frameworkManager := framework_manager.NewFrameworkManager(fakeMasterUpid)
	fakeTpi := &layerx_tpi_client.LayerXTpi{
		CoreURL: "127.0.0.1:34443",
	}
	masterServer := NewMesosApiServerWrapper(fakeTpi, frameworkManager)

	statuses := []*mesosproto.TaskStatus{
		fakes.FakeTaskStatus("task_id_1", mesosproto.TaskState_TASK_FINISHED),
		fakes.FakeTaskStatus("task_id_2", mesosproto.TaskState_TASK_FAILED),
		fakes.FakeTaskStatus("task_id_3", mesosproto.TaskState_TASK_ERROR),
	}

	driverErrc := make(chan error)

	m := masterServer.WrapWithMesos(lxmartini.QuietMartini(), "master@127.0.0.1:3031", driverErrc)
	go m.RunOnAddr(fmt.Sprintf(":3031"))
	go fakes.RunFakeFrameworkServer("fakeframework", 3001)
	go core_fakes.RunFakeLayerXServer(statuses, 34443)
	go func(){
		for {
			err := <- driverErrc
			if err != nil {
				logrus.WithFields(logrus.Fields{"err": err}).Errorf("SHOULD BE TESTING THIS ERROR!")
			}
		}
	}()
	logrus.SetLevel(logrus.DebugLevel)

	Describe("GET " + GET_MASTER_STATE, func() {
		It("returns state of the faux master", func() {
			resp, data, err := lxhttpclient.Get("127.0.0.1:3031", GET_MASTER_STATE, nil)
			Expect(err).To(BeNil())
			var state mesos_data.MesosState
			Expect(resp.StatusCode).To(Equal(200))
			err = json.Unmarshal(data, &state)
			Expect(err).To(BeNil())
			Expect(state.Version).To(Equal("0.25.0"))
			Expect(state.Leader).To(Equal("master@127.0.0.1:3031"))
		})
	})
	Describe("GET " + GET_MASTER_STATE_DEPRECATED, func() {
		It("returns state of the faux master", func() {
			resp, data, err := lxhttpclient.Get("127.0.0.1:3031", GET_MASTER_STATE_DEPRECATED, nil)
			Expect(err).To(BeNil())
			var state mesos_data.MesosState
			Expect(resp.StatusCode).To(Equal(200))
			err = json.Unmarshal(data, &state)
			Expect(err).To(BeNil())
			Expect(state.Version).To(Equal("0.25.0"))
			Expect(state.Leader).To(Equal("master@127.0.0.1:3031"))
		})
	})
	Describe("POST {subscribe_call} " + MESOS_SCHEDULER_CALL, func() {
		It("registers the framework to layer-x core, returns \"FrameworkRegisteredMessage\" to framework", func() {
			fakeSubscribe := fakes.FakeSubscribeCall()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", MESOS_SCHEDULER_CALL, headers, fakeSubscribe)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST {decline_call} " + MESOS_SCHEDULER_CALL, func() {
		It("declines the task-collection offer(s) returns 202 to framework", func() {
			fakeDecline := fakes.FakeDeclineOffersCall("fakeframework", "fake_offer_id")
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", MESOS_SCHEDULER_CALL, headers, fakeDecline)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST {reconcile_tasks} " + MESOS_SCHEDULER_CALL, func() {
		It("tells the layer-x core to bubble up status updates", func() {
			fakeRegisterRequest := fakes.FakeRegisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", REGISTER_FRAMEWORK_MESSAGE, headers, fakeRegisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
			reconcile := fakes.FakeReconcileTasksCall("fakeframework", "task_id_1", "task_id_2", "task_id_3")
			resp, _, err = lxhttpclient.Post("127.0.0.1:3031", MESOS_SCHEDULER_CALL, headers, reconcile)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST {REVIVE_OFFERS} " + MESOS_SCHEDULER_CALL, func() {
		It("does a no op (we are re-sending offers anyway)", func() {
			fakeRegisterRequest := fakes.FakeRegisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", REGISTER_FRAMEWORK_MESSAGE, headers, fakeRegisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
			revive := fakes.FakeReviveOffersCall("fakeframework")
			resp, _, err = lxhttpclient.Post("127.0.0.1:3031", MESOS_SCHEDULER_CALL, headers, revive)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST {ACCEPT_OFFERS} " + MESOS_SCHEDULER_CALL, func() {
		It("submits tasks to layerx core", func() {
			fakeRegisterRequest := fakes.FakeRegisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", REGISTER_FRAMEWORK_MESSAGE, headers, fakeRegisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
			fakeFrameworkId := fakeRegisterRequest.GetFramework().GetId().GetValue()
			fakeTask1 := fakes.FakeMesosTask("fake_task_1")
			fakeTask2 := fakes.FakeMesosTask("fake_task_2")
			fakeTask3 := fakes.FakeMesosTask("fake_task_3")
			fakeLaunchTasks := fakes.FakeLaunchTasksCall(fakeFrameworkId, []string{"fake_offer_id"}, fakeTask1, fakeTask2, fakeTask3)
			resp, _, err = lxhttpclient.Post("127.0.0.1:3031", MESOS_SCHEDULER_CALL, headers, fakeLaunchTasks)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST {ACCEPT_OFFERS} " + MESOS_SCHEDULER_CALL, func() {
		It("sends a kill task request to layer-x", func() {
			fakeRegisterRequest := fakes.FakeRegisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", REGISTER_FRAMEWORK_MESSAGE, headers, fakeRegisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
			fakeFrameworkId := fakeRegisterRequest.GetFramework().GetId().GetValue()
			fakeLaunchTasks := fakes.FakeLaunchTasksMessage(fakeFrameworkId)
			resp, _, err = lxhttpclient.Post("127.0.0.1:3031", LAUNCH_TASKS_MESSAGE, headers, fakeLaunchTasks)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
			fakeTaskId := fakeLaunchTasks.GetTasks()[0].GetTaskId().GetValue()
			fakeKillTask := fakes.FakeKillTaskCall(fakeFrameworkId, fakeTaskId)
			resp, _, err = lxhttpclient.Post("127.0.0.1:3031", MESOS_SCHEDULER_CALL, headers, fakeKillTask)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST " + REGISTER_FRAMEWORK_MESSAGE, func() {
		It("registers the framework to layer-x core, returns \"FrameworkRegisteredMessage\" to framework", func() {
			fakeRegisterRequest := fakes.FakeRegisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", REGISTER_FRAMEWORK_MESSAGE, headers, fakeRegisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST " + REREGISTER_FRAMEWORK_MESSAGE, func() {
		It("registers the framework to layer-x core, returns \"FrameworkRegisteredMessage\" to framework", func() {
			fakeReregisterRequest := fakes.FakeReregisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", REREGISTER_FRAMEWORK_MESSAGE, headers, fakeReregisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST " + UNREGISTER_FRAMEWORK_MESSAGE, func() {
		It("signals layer-x to delete the framework", func() {
			fakeUnregisterRequest := fakes.FakeUnregisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", UNREGISTER_FRAMEWORK_MESSAGE, headers, fakeUnregisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST " + LAUNCH_TASKS_MESSAGE, func() {
		It("submits tasks to layerx core", func() {
			fakeRegisterRequest := fakes.FakeRegisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", REGISTER_FRAMEWORK_MESSAGE, headers, fakeRegisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
			fakeFrameworkId := fakeRegisterRequest.GetFramework().GetId().GetValue()
			fakeLaunchTasks := fakes.FakeLaunchTasksMessage(fakeFrameworkId)
			resp, _, err = lxhttpclient.Post("127.0.0.1:3031", LAUNCH_TASKS_MESSAGE, headers, fakeLaunchTasks)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST " + RECONCILE_TASKS_MESSAGE, func() {
		It("tells the layer-x core to bubble up status updates", func() {
			fakeRegisterRequest := fakes.FakeRegisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", REGISTER_FRAMEWORK_MESSAGE, headers, fakeRegisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
			fakeFrameworkId := fakeRegisterRequest.GetFramework().GetId().GetValue()
			statuses := []*mesosproto.TaskStatus{
				fakes.FakeTaskStatus("task_id_1", mesosproto.TaskState_TASK_RUNNING),
				fakes.FakeTaskStatus("task_id_2", mesosproto.TaskState_TASK_RUNNING),
				fakes.FakeTaskStatus("task_id_3", mesosproto.TaskState_TASK_RUNNING),
			}
			fakeReconcile := fakes.FakeReconcileTasksMessage(fakeFrameworkId, statuses)
			resp, _, err = lxhttpclient.Post("127.0.0.1:3031", RECONCILE_TASKS_MESSAGE, headers, fakeReconcile)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST " + KILL_TASK_MESSAGE, func() {
		It("sends a kill task request to layer-x", func() {
			fakeRegisterRequest := fakes.FakeRegisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", REGISTER_FRAMEWORK_MESSAGE, headers, fakeRegisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
			fakeFrameworkId := fakeRegisterRequest.GetFramework().GetId().GetValue()
			fakeLaunchTasks := fakes.FakeLaunchTasksMessage(fakeFrameworkId)
			resp, _, err = lxhttpclient.Post("127.0.0.1:3031", LAUNCH_TASKS_MESSAGE, headers, fakeLaunchTasks)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
			fakeTaskId := fakeLaunchTasks.GetTasks()[0].GetTaskId().GetValue()
			fakeKillTaskMessage := fakes.FakeKillTaskMessage(fakeFrameworkId, fakeTaskId)
			resp, _, err = lxhttpclient.Post("127.0.0.1:3031", KILL_TASK_MESSAGE, headers, fakeKillTaskMessage)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST " + STATUS_UPDATE_ACKNOWLEDGEMENT_MESSAGE, func() {
		It("logs the request to debug (noop)", func() {
			fakeStatusUpdateAck := fakes.FakeStatusUpdateAcknowledgementMessage("doesnt_matter_fwid", "doesntmattertaskid", "any_slave", []byte("some_bytes"))
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", STATUS_UPDATE_ACKNOWLEDGEMENT_MESSAGE, headers, fakeStatusUpdateAck)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
	Describe("POST " + REVIVE_OFFERS_MESSAGE, func() {
		It("logs the request to debug (noop)", func() {
			fakeReviveOffersMsg := fakes.FakeReviveOffersMessage("doesnt_matter_fwid")
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", REVIVE_OFFERS_MESSAGE, headers, fakeReviveOffersMsg)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
})
