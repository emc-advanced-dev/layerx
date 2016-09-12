package kube

import (
	"k8s.io/client-go/1.4/kubernetes"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/pkg/api"
	"github.com/emc-advanced-dev/pkg/errors"
)

type Client struct {
	kubeClient *kubernetes.Clientset
}

func NewClient(kubeClient *kubernetes.Clientset) *Client{
	return &Client{
		kubeClient: kubeClient,
	}
}

func (c *Client) FetchNodes() ([]*lxtypes.Node, error) {
	nodes := []*lxtypes.Node{}
	//collect available kubeletes
	nl, err := c.kubeClient.Core().Nodes().List(api.ListOptions{})
	if err != nil {
		return nil, errors.New("getting kube nodes list", err)
	}
	for _, n := range nl.Items {
		nodeId := string(n.UID)
		kCpus := n.Status.Allocatable[v1.ResourceCPU]
		cpus := float64((&kCpus).Value())
		kMemMB := n.Status.Allocatable[v1.ResourceMemory]
		memMB := float64((&kMemMB).Value() >> 10)
		kDiskMB := n.Status.Allocatable[v1.ResourceStorage]
		diskMB := float64((&kDiskMB).Value() >> 10)

		resource := &lxtypes.Resource{
			Id: nodeId,
			NodeId: nodeId,
			Cpus: cpus,
			Mem: memMB,
			Disk: diskMB,
			//TODO: figure out how to get port resource from k8s nodes
			Ports: []lxtypes.PortRange{
				lxtypes.PortRange{
					Begin: 1,
					End: 65535,
				},
			},
		}

		lxNode := lxtypes.NewNode(nodeId)
		lxNode.AddResource(resource)
		nodes = append(nodes, lxNode)
	}

	return nodes, nil
}

func (c *Client) LaunchTasks(launchTasksMessage layerx_rpi_client.LaunchTasksMessage) error {
	return nil
}

func (c *Client) KillTask(taskID string) error {
	return nil
}