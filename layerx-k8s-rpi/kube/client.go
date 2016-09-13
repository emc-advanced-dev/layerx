package kube

import (
	"k8s.io/client-go/1.4/kubernetes"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/pkg/api"
	"github.com/emc-advanced-dev/pkg/errors"
	"fmt"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/Sirupsen/logrus"
	"k8s.io/client-go/1.4/pkg/api/resource"
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
	if len(launchTasksMessage.ResourcesToUse) < 0 {
		return errors.New("must specify a node to launch these tasks on", nil)
	}
	nodeName := launchTasksMessage.ResourcesToUse[0].Id
	if len(launchTasksMessage.ResourcesToUse) > 1 {
		logrus.Warn("K8s RPI only supports launching tasks on a single node at a time. using node", nodeName)
	}
	for _, task := range launchTasksMessage.TasksToLaunch {
		pod, err := convertToPod(task, nodeName)
		if err != nil {
			return errors.New("failed to convert task to pod", err)
		}
		logrus.Debug("creating pod", pod)
		result, err := c.kubeClient.Core().Pods("").Create(pod)
		if err != nil {
			return errors.New("failed to create pod on k8s", err)
		}
		logrus.Infof("created pod", result)
	}

	return nil
}

func (c *Client) KillTask(taskID string) error {
	return nil
}

func convertToPod(task *lxtypes.Task, nodeName string) (*v1.Pod, error) {
	logrus.Debug("converting task", task, "to pod")
	if task.Container == nil || task.Container.Docker == nil {
		return nil, errors.New("only tasks with docker images are supported for the k8s rpi", nil)
	}

	objectMeta := v1.ObjectMeta{
		Name: task.TaskId,
	}
	kubePorts := []v1.ContainerPort{}
	for _, port := range task.Container.Docker.PortMappings {
		protocol := v1.ProtocolTCP
		if *port.Protocol == "udp" {
			protocol = v1.ProtocolUDP
		}
		kubePort := v1.ContainerPort{
			ContainerPort: int32(*port.ContainerPort),
			HostPort: int32(*port.HostPort),
			Protocol: protocol,
		}
		kubePorts = append(kubePorts, kubePort)
	}
	kubeEnvVars := []v1.EnvVar{}
	for _, env := range task.Command.Environment.Variables {
		kubeVar := v1.EnvVar{
			Name: *env.Name,
			Value: *env.Value,
		}
		kubeEnvVars = append(kubeEnvVars, kubeVar)
	}

	cpus, err := resource.ParseQuantity(fmt.Sprintf("%v", task.Cpus))
	if err != nil {
		return nil, errors.New("parsing quantity from task", err)
	}
	mem, err := resource.ParseQuantity(fmt.Sprintf("%v", task.Mem * 1024))
	if err != nil {
		return nil, errors.New("parsing quantity from task", err)
	}
	disk, err := resource.ParseQuantity(fmt.Sprintf("%v", task.Disk * 1024))
	if err != nil {
		return nil, errors.New("parsing quantity from task", err)
	}
	resources := v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU: cpus,
			v1.ResourceMemory: mem,
			v1.ResourceStorage: disk,
		},
	}

	kubeMounts := []v1.VolumeMount{}
	kubeVols := []v1.Volume{}
	for i, volume := range task.Container.Volumes {
		readOnly := true
		if *volume.Mode == mesosproto.Volume_RW {
			readOnly = false
		}
		volName := fmt.Sprintf("%s-%d", task.Name, i)
		kubeVol := v1.Volume{
			Name: volName,
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: *volume.HostPath,
				},
			},
		}
		kubeMount := v1.VolumeMount{
			Name: volName,
			ReadOnly: readOnly,
			MountPath: *volume.ContainerPath,
		}
		kubeVols = append(kubeVols, kubeVol)
		kubeMounts = append(kubeMounts, kubeMount)
	}

	container := v1.Container{
		Name: task.Name,
		Image: *task.Container.Docker.Image,
		Args: task.Command.Arguments,
		Ports: kubePorts,
		Env: kubeEnvVars,
		Resources: resources,
		VolumeMounts: kubeMounts,
	}
	spec := v1.PodSpec{
		Volumes: kubeVols,
		Containers: []v1.Container{container},
		RestartPolicy: v1.RestartPolicyNever,
		NodeName: nodeName,
	}
	return &v1.Pod{
		ObjectMeta: objectMeta,
		Spec: spec,
	}, nil
}