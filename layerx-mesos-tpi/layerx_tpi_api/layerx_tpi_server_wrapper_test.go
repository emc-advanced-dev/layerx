package layerx_tpi_api_test

import (
	. "github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/layerx_tpi_api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/fakes"
	"github.com/Sirupsen/logrus"
	"fmt"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/framework_manager"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/mesos_master_api/mesos_data"
	core_fakes "github.com/emc-advanced-dev/layerx/layerx-core/fakes"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/mesos_master_api"
	"github.com/mesos/mesos-go/mesosproto"
	"net/http"
)

var _ = Describe("LayerxTpiServerWrapper", func() {
	fakeMasterUpid, _ := mesos_data.UPIDFromString("master@127.0.0.1:3032")
	frameworkManager := framework_manager.NewFrameworkManager(fakeMasterUpid)
	fakeTpi := &layerx_tpi_client.LayerXTpi{
		CoreURL: "127.0.0.1:34445",
	}
	tpiServerWrapper := NewTpiApiServerWrapper(fakeTpi, frameworkManager)

	masterServer := mesos_master_api.NewMesosApiServerWrapper(fakeTpi, frameworkManager)

	m := tpiServerWrapper.WrapWithTpi(lxmartini.QuietMartini(), "master@127.0.0.1:3032", make(chan error))
	m = masterServer.WrapWithMesos(m, "master@127.0.0.1:3032", make(chan error))
	go m.RunOnAddr(fmt.Sprintf(":3032"))
	go fakes.RunFakeFrameworkServer("fakeframework", 3002)
	go core_fakes.RunFakeLayerXServer(nil, 34445)
	logrus.SetLevel(logrus.DebugLevel)

	Describe("POST {collect_tasks_message} " + COLLECT_TASKS, func() {
		It("sends collect_task_message to the framework", func() {
			fakeRegisterRequest := fakes.FakeRegisterFrameworkMessage()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3002",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3032", mesos_master_api.REGISTER_FRAMEWORK_MESSAGE, headers, fakeRegisterRequest)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))

			fakeCollectTasksMsg := &layerx_tpi_client.CollectTasksMessage{
				TaskProviders: []*lxtypes.TaskProvider{
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
			fakeUpdateTaskStatusMessage := &layerx_tpi_client.UpdateTaskStatusMessage{
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

	Describe("POST {HealthCheckTaskProviderMessage} " + HEALTH_CHECK_TASK_PROVIDER, func() {
		Context("the framework is not connected", func(){
			It("performs a health check on the target framework and responds with 410", func() {
				fakeHealthCheckTaskProviderMessage := &layerx_tpi_client.HealthCheckTaskProviderMessage{
					TaskProvider: &lxtypes.TaskProvider{
						Id: "fakedisconnectedframework",
						Source: "fakedisconnectedframework@127.0.0.1:1987",
					},
				}
				resp, _, err := lxhttpclient.Post("127.0.0.1:3032", HEALTH_CHECK_TASK_PROVIDER, nil, fakeHealthCheckTaskProviderMessage)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusGone))
			})
		})
		Context("the framework is connected", func(){
			It("performs a health check on the target framework and responds with 200", func() {
				fakeHealthCheckTaskProviderMessage := &layerx_tpi_client.HealthCheckTaskProviderMessage{
					TaskProvider: &lxtypes.TaskProvider{
						Id: "fake_task_provider_id",
						Source: "fakeframework@127.0.0.1:3002",
					},
				}
				resp, _, err := lxhttpclient.Post("127.0.0.1:3032", HEALTH_CHECK_TASK_PROVIDER, nil, fakeHealthCheckTaskProviderMessage)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

})
