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
	"strings"
	"github.com/golang/protobuf/proto"
	"time"
)

var (
	watchTimeout = time.Minute * 3
)

const (
	namespaceName = "layer-x"
)

type Client struct {
	kubeClient *kubernetes.Clientset
}

func NewClient(kubeClient *kubernetes.Clientset) *Client{
	return &Client{
		kubeClient: kubeClient,
	}
}

func (c *Client) Init() error {
	//create namespace
	namespace := &v1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: namespaceName,
		},
	}
	logrus.Debug("creating namespace", namespace)
	res, err := c.kubeClient.Core().Namespaces().Create(namespace)
	if err != nil {
		return errors.New("creating namespace", err)
	}
	logrus.Debug("successfully created", res)
	return nil
}

func (c *Client) Teardown() error {
	//delete namespace
	logrus.Debug("deleting namespace", namespaceName)
	if err := c.kubeClient.Core().Namespaces().Delete(namespaceName, &api.DeleteOptions{}); err != nil {
		return errors.New("deleting namespace", err)
	}
	logrus.Debug("successfully deleted namespace", namespaceName)
	return nil
}

func (c *Client) FetchNodes() ([]*lxtypes.Node, error) {
	nodes := []*lxtypes.Node{}
	//collect available kubeletes
	nl, err := c.kubeClient.Core().Nodes().List(api.ListOptions{})
	if err != nil {
		return nil, errors.New("getting kube nodes list", err)
	}
	for _, n := range nl.Items {
		nodeId := n.Name
		kCpus := n.Status.Allocatable[v1.ResourceCPU]
		cpus := float64((&kCpus).Value())
		kMemMB := n.Status.Allocatable[v1.ResourceMemory]
		memMB := float64((&kMemMB).Value())
		kDiskMB := n.Status.Allocatable[v1.ResourceStorage]
		diskMB := float64((&kDiskMB).Value())

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
	if len(launchTasksMessage.ResourcesToUse) < 1 {
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
		result, err := c.kubeClient.Core().Pods(namespaceName).Create(pod)
		if err != nil {
			return errors.New("failed to create pod on k8s", err)
		}
		if err := c.waitPodCreate(result.Name); err != nil {
			return errors.New("waiting for pod to be created", err)
		}
		logrus.Infof("created pod", result)
	}

	return nil
}

func (c *Client) KillTask(taskID string) error {
	logrus.Debug("killing task", taskID)
	if err := c.kubeClient.Core().Pods(namespaceName).Delete(taskID, &api.DeleteOptions{}); err != nil {
		return errors.New("deleting pod "+taskID, err)
	}
	if err := c.waitPodDelete(taskID); err != nil {
		return errors.New("waiting for pod to be deleted", err)
	}
	return nil
}

func (c *Client) GetStatuses() ([]*mesosproto.TaskStatus, error) {
	logrus.Debug("getting status for all pods")
	podList, err := c.kubeClient.Core().Pods(namespaceName).List(api.ListOptions{})
	if err != nil {
		return nil, errors.New("getting pod list", err)
	}
	statuses := []*mesosproto.TaskStatus{}
	for _, pod := range podList.Items {
		statuses = append(statuses, getStatus(pod))
	}
	logrus.Debug("got status list", statuses)
	return statuses, nil
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
	if task.Command.Environment != nil {
		for _, env := range task.Command.Environment.Variables {
			kubeVar := v1.EnvVar{
				Name: *env.Name,
				Value: *env.Value,
			}
			kubeEnvVars = append(kubeEnvVars, kubeVar)
		}
	}
	cpus, err := resource.ParseQuantity(fmt.Sprintf("%v", task.Cpus))
	if err != nil {
		return nil, errors.New("parsing quantity from task", err)
	}

	mem, err := resource.ParseQuantity(fmt.Sprintf("%vMi", task.Mem))
	if err != nil {
		return nil, errors.New("parsing quantity from task", err)
	}
	//disk, err := resource.ParseQuantity(fmt.Sprintf("%v", task.Disk * 1024))
	//if err != nil {
	//	return nil, errors.New("parsing quantity from task", err)
	//}
	resources := v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU: cpus,
			v1.ResourceMemory: mem,
			//todo: figure out if it's possible to set storage limits for containers
			//v1.ResourceStorage: disk,
		},
	}

	logrus.Debugf("task %v has limits %v", task, resources)

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

	//args or command for container
	args := []string{}
	if binary := *task.Command.Value; binary != "" {
		args = append(args, strings.Split(binary, " ")...)
	}
	args = append(args, task.Command.Arguments...)
	cmd := []string{}
	// no entrypoint
	if *task.Command.Shell {
		cmd = args
		args = []string{}
	}

	container := v1.Container{
		Name: task.Name,
		Image: *task.Container.Docker.Image,
		Command: cmd,
		Args: args,
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

func getStatus(pod v1.Pod) *mesosproto.TaskStatus {
	name := pod.Name
	message := pod.Status.Message
	var mesosState mesosproto.TaskState
	switch pod.Status.Phase {
	case v1.PodPending:
		mesosState = mesosproto.TaskState_TASK_STAGING
	case v1.PodRunning:
		mesosState = mesosproto.TaskState_TASK_RUNNING
	case v1.PodSucceeded:
		mesosState = mesosproto.TaskState_TASK_FINISHED
	case v1.PodFailed:
		mesosState = mesosproto.TaskState_TASK_FAILED
	case v1.PodUnknown:
		mesosState = mesosproto.TaskState_TASK_ERROR
	}
	nodeID := pod.Spec.NodeName
	return &mesosproto.TaskStatus{
		TaskId: &mesosproto.TaskID{Value: proto.String(name)},
		State: &mesosState,
		Message: proto.String(message),
		SlaveId: &mesosproto.SlaveID{Value: proto.String(nodeID)},
	}
}

func (c *Client) waitPodCreate(name string) error {
	return PollWait(func() bool{
		_, err := c.kubeClient.Core().Pods(namespaceName).Get(name)
		if err == nil {
			return true
		}
		return false
	})
}

func (c *Client) waitPodDelete(name string) error {
	return PollWait(func() bool{
			_, err := c.kubeClient.Core().Pods(namespaceName).Get(name)
			if err != nil {
				return true
			}
		return false
	})
}

func PollWait(waitFunc func() bool) error {
	finished := make(chan struct{})
	go func(){
		for {
			if waitFunc() {
				close(finished)
				return
			}
			time.Sleep(time.Second)
		}
	}()
	select {
	case <-finished:
		return nil
	case <-time.After(watchTimeout):
		return errors.New("waiting for result timed out", nil)
	}
}