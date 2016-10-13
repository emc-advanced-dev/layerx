package swarm_test

import (
	. "github.com/emc-advanced-dev/layerx/layerx-swarm-rpi/swarm"

	"fmt"
	"github.com/Sirupsen/logrus"
	core_fakes "github.com/emc-advanced-dev/layerx/layerx-core/fakes"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/mesos/mesos-go/mesosproto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strings"
	"github.com/emc-advanced-dev/layerx/layerx-k8s-rpi/kube"
)

var (
	client       *Client
	fakeCore     *core_fakes.FakeCore
	started      bool
	fakeCorePort = 6123
)

var _ = Describe("Client", func() {
	logrus.SetLevel(logrus.DebugLevel)
	BeforeEach(func() {
		if !started {
			fakeCore = core_fakes.NewFakeCore()
			go fakeCore.Start(nil, fakeCorePort)
			started = true
			if err := setUp(); err != nil && !strings.Contains(err.Error(), "already exists") {
				fmt.Println(err)
				os.Exit(-1)
			}
		}
	})
	Describe("FetchResources", func() {
		It("returns all the kube nodes as an array of lx resourecs", func() {
			resources, err := client.FetchResources()
			Expect(err).To(BeNil())
			Expect(resources).NotTo(BeEmpty())
			fmt.Printf("Nodes: %+v", resources[0])
		})
	})
	Describe("LaunchTasks", func() {
		It("launches lx task as a pod on the target k8s node", func() {
			nodes, err := client.FetchResources()
			Expect(err).To(BeNil())
			Expect(nodes).NotTo(BeEmpty())
			fakeTask := core_fakes.FakeLXDockerTask("id-1234", "fake-task", nodes[0].Id, "echo DID IT WORKED??")
			fakeTask.Mem = 4
			launchTasksMessage := layerx_rpi_client.LaunchTasksMessage{
				TasksToLaunch:  []*lxtypes.Task{fakeTask},
				ResourcesToUse: []*lxtypes.Resource{nodes[0]},
			}
			err = client.LaunchTasks(launchTasksMessage)
			defer client.KillTask(fakeTask.TaskId)
			Expect(err).To(BeNil())
		})
	})
	Describe("GetStatuses", func() {
		It("gets the status for all existing pods", func() {
			nodes, err := client.FetchResources()
			Expect(err).To(BeNil())
			Expect(nodes).NotTo(BeEmpty())
			fakeTask := core_fakes.FakeLXDockerTask("id-1234", "fake-task", nodes[0].Id, "echo STARTING! && sleep 1 && echo FINISHED!")
			fakeTask.Mem = 4
			launchTasksMessage := layerx_rpi_client.LaunchTasksMessage{
				TasksToLaunch:  []*lxtypes.Task{fakeTask},
				ResourcesToUse: []*lxtypes.Resource{nodes[0]},
			}
			err = client.LaunchTasks(launchTasksMessage)
			defer client.KillTask(fakeTask.TaskId)
			Expect(err).To(BeNil())
			statuses, err := client.GetStatuses()
			Expect(err).To(BeNil())
			Expect(statuses).NotTo(BeEmpty())
			Expect(*statuses[0].State).To(Equal(mesosproto.TaskState_TASK_STARTING))
			statuses, err = client.GetStatuses()
			Expect(err).To(BeNil())
			Expect(statuses).NotTo(BeEmpty())
			err = kube.PollWait(func() bool {
				state, err := getFirstTaskState()
				if err != nil {
					logrus.Error("failed to get first task status", err)
					return false
				}
				return *state == mesosproto.TaskState_TASK_FINISHED
			})
			Expect(err).To(BeNil())
		})
	})
	Describe("KillTask", func() {
		It("Calls CoreMessenger.SubmitResource() with an array of lx resourecs", func() {
			nodes, err := client.FetchResources()
			Expect(err).To(BeNil())
			Expect(nodes).NotTo(BeEmpty())
			fakeTask := core_fakes.FakeLXDockerTask("id-1234", "fake-task", nodes[0].Id, "echo DID IT WORKED??")
			fakeTask.Mem = 4
			launchTasksMessage := layerx_rpi_client.LaunchTasksMessage{
				TasksToLaunch:  []*lxtypes.Task{fakeTask},
				ResourcesToUse: []*lxtypes.Resource{nodes[0]},
			}
			err = client.LaunchTasks(launchTasksMessage)
			defer client.KillTask(fakeTask.TaskId)
			Expect(err).To(BeNil())
			err = client.KillTask(fakeTask.TaskId)
			Expect(err).To(BeNil())
		})
	})
})

func getFirstTaskState() (*mesosproto.TaskState, error) {
	statuses, err := client.GetStatuses()
	if err != nil {
		return nil, errors.New("error polling statuses", err)
	}
	if len(statuses) < 1 {
		return nil, errors.New("status length < 1", nil)
	}
	return statuses[0].State, nil
}

func setUp() (err error) {
	client, err = NewClient("fake-swarm-rpi")
	return
}
