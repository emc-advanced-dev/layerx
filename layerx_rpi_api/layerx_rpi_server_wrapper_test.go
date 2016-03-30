package layerx_rpi_api_test

import (
	. "github.com/layer-x/layerx-mesos-rpi_v2/layerx_rpi_api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-commons/lxmartini"
	"fmt"
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/gogo/protobuf/proto"
	core_fakes "github.com/layer-x/layerx-core_v2/fakes"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/layer-x/layerx-mesos-rpi_v2/mesos_framework_api"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/mesos/mesos-go/scheduler"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/Sirupsen/logrus"
	"os"
	"github.com/layer-x/layerx-core_v2/lxtypes"
	"github.com/layer-x/layerx-mesos-rpi_v2/fakes"
	"time"
)

var _ = Describe("LayerxRpiServerWrapper", func() {

	fakeRpi := &layerx_rpi_client.LayerXRpi{
		CoreURL: "127.0.0.1:34446",
	}

	Describe("Setup", func(){
		It("sets up for the test", func(){
			mesosUrl := "127.0.0.1:5050"
			if os.Getenv("MESOS_URL") != "" {
				mesosUrl = os.Getenv("MESOS_URL")
			}
			fakeFramework := &mesosproto.FrameworkInfo{
				User: proto.String(""),
				FailoverTimeout: proto.Float64(15),
				Name: proto.String("FAKE Layer-X Mesos RPI Framework"),
			}
			fakeRpiScheduler := mesos_framework_api.NewRpiMesosScheduler(fakeRpi)
			config := scheduler.DriverConfig{
				Scheduler:  fakeRpiScheduler,
				Framework:  fakeFramework,
				Master:     mesosUrl,
				HostnameOverride: "localhost",
				Credential: (*mesosproto.Credential)(nil),
			}

			driver, err := scheduler.NewMesosSchedulerDriver(config)
			if err != nil {
				err = lxerrors.New("initializing mesos schedulerdriver", err)
				lxlog.Errorf(logrus.Fields{
					"error":     err,
					"mesos_url": mesosUrl,
				}, "error initializing mesos schedulerdriver")
			}
			Expect(err).To(BeNil())

			go func() {
				status, err := driver.Run()
				if err != nil {
					err = lxerrors.New("Framework stopped with status " + status.String(), err)
					lxlog.Errorf(logrus.Fields{
						"error":     err,
						"mesos_url": mesosUrl,
					}, "error running mesos schedulerdriver")
					panic(err)
				}
			}()
			mesosSchedulerDriver := fakeRpiScheduler.GetDriver()
			rpiServerWrapper := NewRpiApiServerWrapper(fakeRpi, mesosSchedulerDriver)

			m := rpiServerWrapper.WrapWithRpi(lxmartini.QuietMartini(), make(chan error))
			go core_fakes.RunFakeLayerXServer(nil, 34446)
			go m.RunOnAddr(fmt.Sprintf(":3033"))
			lxlog.ActiveDebugMode()
		})
	})
	Describe("POST " + COLLECT_RESOURCES, func() {
		It("tells the rpi to send ReviveOffers() to mesos", func() {
			resp, _, err := lxhttpclient.Post("127.0.0.1:3033", COLLECT_RESOURCES, nil, nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})

	Describe("POST {launch_tasks_message} " + LAUNCH_TASKS, func() {
		It("launches tasks on mesos", func() {
			nodes, err := fakeRpi.GetNodes()
			Expect(err).To(BeNil())
			realResources := nodes[0].GetResources()
			fakeTask1 := fakes.FakeTask("fake_task_id_1")
			fakeTask2 := fakes.FakeTask("fake_task_id_2")
			fakeTask3 := fakes.FakeTask("fake_task_id_3")
			fakeTask1.NodeId = realResources[0].NodeId
			fakeTask2.NodeId = realResources[0].NodeId
			fakeTask3.NodeId = realResources[0].NodeId
			fakeTasks := []*lxtypes.Task{fakeTask1, fakeTask2, fakeTask3}
			launchTasksMessage := layerx_rpi_client.LaunchTasksMessage{
				TasksToLaunch: fakeTasks,
				ResourcesToUse: realResources,
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3033", LAUNCH_TASKS, nil, launchTasksMessage)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
			time.Sleep(3000 * time.Millisecond)
		})
	})

	Describe("POST {collect_tasks_message} " + LAUNCH_TASKS, func() {
		It("sends collect_task_message to the framework", func() {
			resp, _, err := lxhttpclient.Post("127.0.0.1:3033", COLLECT_RESOURCES, nil, nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
})
