package layerx_tpi_api_test

import (
	. "github.com/layer-x/layerx-mesos-tpi_v2/layerx_tpi_api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-mesos-tpi_v2/fakes"
	"github.com/layer-x/layerx-commons/lxlog"
	"fmt"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-mesos-tpi_v2/driver"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	core_fakes "github.com/layer-x/layerx-core_v2/fakes"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api"
	"github.com/mesos/mesos-go/mesosproto"
)

var _ = Describe("LayerxTpiServerWrapper", func() {
	actionQueue := lxactionqueue.NewActionQueue()
	fakeMasterUpid, _ := mesos_data.UPIDFromString("master@127.0.0.1:3032")
	frameworkManager := framework_manager.NewFrameworkManager(fakeMasterUpid)
	fakeTpi := &layerx_tpi.LayerXTpi{
		CoreURL: "127.0.0.1:34445",
	}
	tpiServerWrapper := NewTpiApiServerWrapper(fakeTpi, actionQueue, frameworkManager)
	driver := driver.NewMesosTpiDriver(actionQueue)

	masterServer := mesos_master_api.NewMesosApiServerWrapper(fakeTpi, actionQueue, frameworkManager)

	m := tpiServerWrapper.WrapWithTpi(lxmartini.QuietMartini(), "master@127.0.0.1:3032", make(chan error))
	m = masterServer.WrapWithMesos(m, "master@127.0.0.1:3032", make(chan error))
	go m.RunOnAddr(fmt.Sprintf(":3032"))
	go driver.Run()
	go fakes.RunFakeFrameworkServer("fakeframework", 3002)
	go core_fakes.RunFakeLayerXServer(nil, 34445)
	lxlog.ActiveDebugMode()

	Describe("POST {collect_tasks_message} " + COLLECT_TASKS, func() {
		It("sends collect_task_message to the framework", func() {
			fakeRegisterRequest := fakes.FakeRegisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3002",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3032", mesos_master_api.REGISTER_FRAMEWORK_MESSAGE, headers, fakeRegisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))

			fakeCollectTasksMsg := &layerx_tpi.CollectTasksMessage{
				TaskProivders: []*lxtypes.TaskProvider{
					&lxtypes.TaskProvider{
						Id: "fake_task_provider_id",
						Source: "fakeframework@127.0.0.1:3002",
					},
				},
			}
			resp, _, err = lxhttpclient.Post("127.0.0.1:3032", COLLECT_TASKS, nil, fakeCollectTasksMsg)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})

	Describe("POST {UpdateTaskStatusMessage} " + UPDATE_TASK_STATUS, func() {
		It("sends status update to the framework", func() {
			fakeUpdateTaskStatusMessage := &layerx_tpi.UpdateTaskStatusMessage{
				TaskProvider: &lxtypes.TaskProvider{
					Id: "fake_task_provider_id",
					Source: "fakeframework@127.0.0.1:3002",
				},
				TaskStatus: fakes.FakeTaskStatus("fake_task_1", mesosproto.TaskState_TASK_RUNNING),
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3032", UPDATE_TASK_STATUS, nil, fakeUpdateTaskStatusMessage)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})

})
