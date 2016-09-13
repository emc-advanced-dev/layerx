package kube_test

import (
	. "github.com/emc-advanced-dev/layerx/layerx-k8s-rpi/kube"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/1.4/tools/clientcmd"
	"k8s.io/client-go/1.4/kubernetes"
	"os"
	"github.com/emc-advanced-dev/pkg/errors"
	core_fakes "github.com/emc-advanced-dev/layerx/layerx-core/fakes"
	"fmt"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
)

var (
	client *Client
	fakeCore *core_fakes.FakeCore
	started bool
	fakeCorePort = 6123
)

var _ = Describe("Client", func() {
	BeforeEach(func() {
		if err := setUp(); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fakeCore = core_fakes.NewFakeCore()
		if !started {
			go fakeCore.Start(nil, fakeCorePort)
			started = true
		}
	})
	Describe("FetchResources", func() {
		It("Calls CoreMessenger.SubmitResource() with an array of lx resourecs", func() {
			nodes, err := client.FetchNodes()
			Expect(err).To(BeNil())
			Expect(nodes).NotTo(BeEmpty())
			fmt.Printf("Nodes: %+v", nodes[0])
		})
	})
	Describe("LaunchTasks", func() {
		It("Calls CoreMessenger.SubmitResource() with an array of lx resourecs", func() {
			launchTasksMessage := layerx_rpi_client.LaunchTasksMessage{
				TasksToLaunch: []*lxtypes.Task{},
				ResourcesToUse: []*lxtypes.Resource{},
			}

			err := client.LaunchTasks(launchTasksMessage)
			Expect(err).To(BeNil())
		})
	})
})

func setUp() error {
	kubeconfig := os.Getenv("KUBE_CFG")
	if kubeconfig == "" {
		return errors.New("path to kubeconfig must be specified by env var KUBE_CFG", nil)
	}
	//initialize kube client
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	client = NewClient(clientset)
	return nil
}
