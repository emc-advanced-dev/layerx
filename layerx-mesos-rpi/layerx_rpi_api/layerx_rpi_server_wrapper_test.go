package layerx_rpi_api_test

import (
	. "github.com/emc-advanced-dev/layerx/layerx-mesos-rpi/layerx_rpi_api"

	"fmt"
	"github.com/Sirupsen/logrus"
	core_fakes "github.com/emc-advanced-dev/layerx/layerx-core/fakes"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-rpi/fakes"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-rpi/mesos_framework_api"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"time"
)

var _ = Describe("LayerxRpiServerWrapper", func() {

	fakeRpi := &layerx_rpi_client.LayerXRpi{
		CoreURL: "127.0.0.1:34446",
	}

	Describe("Setup", func() {
		It("sets up for the test", func() {
			mesosUrl := "127.0.0.1:5050"
			if os.Getenv("MESOS_URL") != "" {
				mesosUrl = os.Getenv("MESOS_URL")
			}
			fakeFramework := &mesosproto.FrameworkInfo{
				User:            proto.String(""),
				FailoverTimeout: proto.Float64(15),
				Name:            proto.String("FAKE Layer-X Mesos RPI Framework"),
			}
			fakeRpiScheduler := mesos_framework_api.NewRpiMesosScheduler(fakeRpi)
			config := scheduler.DriverConfig{
				Scheduler:        fakeRpiScheduler,
				Framework:        fakeFramework,
				Master:           mesosUrl,
				HostnameOverride: "localhost",
				Credential:       (*mesosproto.Credential)(nil),
			}

			driver, err := scheduler.NewMesosSchedulerDriver(config)
			if err != nil {
				err = errors.New("initializing mesos schedulerdriver", err)
				logrus.WithFields(logrus.Fields{
					"error":     err,
					"mesos_url": mesosUrl,
				}).Errorf("error initializing mesos schedulerdriver")
			}
			Expect(err).To(BeNil())

			go func() {
				status, err := driver.Run()
				if err != nil {
					err = errors.New("Framework stopped with status "+status.String(), err)
					logrus.WithFields(logrus.Fields{
						"error":     err,
						"mesos_url": mesosUrl,
					}).Errorf("error running mesos schedulerdriver")
					panic(err)
				}
			}()
			mesosSchedulerDriver := fakeRpiScheduler.GetDriver()
			rpiServerWrapper := NewRpiApiServerWrapper(fakeRpi, mesosSchedulerDriver)

			m := rpiServerWrapper.WrapWithRpi(lxmartini.QuietMartini(), make(chan error))
			go core_fakes.NewFakeCore().Start(nil, 34446)
			go m.RunOnAddr(fmt.Sprintf(":3033"))
			logrus.SetLevel(logrus.DebugLevel)
		})
	})
	Describe("POST "+COLLECT_RESOURCES, func() {
		It("tells the rpi to send ReviveOffers() to mesos", func() {
			resp, _, err := lxhttpclient.Post("127.0.0.1:3033", COLLECT_RESOURCES, nil, nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})

	Describe("POST {launch_tasks_message} "+LAUNCH_TASKS, func() {
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
				TasksToLaunch:  fakeTasks,
				ResourcesToUse: realResources,
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3033", LAUNCH_TASKS, nil, launchTasksMessage)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
			time.Sleep(3000 * time.Millisecond)
		})
	})

	Describe("POST {collect_tasks_message} "+LAUNCH_TASKS, func() {
		It("sends collect_task_message to the framework", func() {
			resp, _, err := lxhttpclient.Post("127.0.0.1:3033", COLLECT_RESOURCES, nil, nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
})
