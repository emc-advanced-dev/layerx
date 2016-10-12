package swarm

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/golang/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
	"strings"
	"time"
	"github.com/fsouza/go-dockerclient"
	"os"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/mount"
)

type Client struct {
	docker  *docker.Client
	rpiName string
}

func NewClient(rpiName string) (*Client, error) {
	host := os.Getenv("DOCKER_HOST")
	path := os.Getenv("DOCKER_CERT_PATH")
	logrus.Infof("using DOCKER_HOST=%v", host)
	logrus.Infof("using DOCKER_CERT_PATH=%v", path)
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)
	client, err := docker.NewTLSClient(host, cert, key, ca)
	if err != nil {
		return nil, errors.New("creating TLS client for docker", err)
	}
	return &Client{
		docker: client,
		rpiName: rpiName,
	}, nil
}

func (c *Client) FetchResources() ([]*lxtypes.Resource, error) {
	nl, err := c.docker.ListNodes(docker.ListNodesOptions{})
	if err != nil {
		return nil, errors.New("getting node list from docker client", err)
	}
	resources := []*lxtypes.Resource{}
	for _, n := range nl {
		nodeId := n.ID
		cpus := float64(n.Description.Resources.NanoCPUs >> 29)
		memMB := float64(n.Description.Resources.MemoryBytes >> 20)
		//todo: figure out how to translate disk on swarm nodes
		diskMB := float64(0)

		resource := &lxtypes.Resource{
			Id:     nodeId,
			NodeId: nodeId,
			Cpus:   cpus,
			Mem:    memMB,
			Disk:   diskMB,
			//TODO: figure out how to get port resource from swarm nodes
			Ports: []lxtypes.PortRange{
				lxtypes.PortRange{
					Begin: 1,
					End:   65535,
				},
			},
			ResourceType: lxtypes.ResourceType_DockerSwarm,
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

func (c *Client) LaunchTasks(launchTasksMessage layerx_rpi_client.LaunchTasksMessage) error {
	c.docker.CreateService()
	if len(launchTasksMessage.ResourcesToUse) < 1 {
		return errors.New("must specify a node to launch these tasks on", nil)
	}
	nodeName := launchTasksMessage.ResourcesToUse[0].Id
	if len(launchTasksMessage.ResourcesToUse) > 1 {
		logrus.Warn("Swarm RPI only supports launching tasks on a single node at a time. using node", nodeName)
	}
	for _, task := range launchTasksMessage.TasksToLaunch {
		if task.Container == nil || task.Container.GetDocker() == nil {
			return errors.New("only tasks with docker images are supported for the docker swarm rpi", nil)
		}
		logrus.Infof("launching task", task)
		svc := c.convertToServiceSpec(task, nodeName)
		res, err := c.docker.CreateService(svc)
		if err != nil {
			return errors.New("launching task", err)
		}
		logrus.Debug("launched service", res)
	}

	return nil
}

func (c *Client) KillTask(taskID string) error {
	logrus.Debug("killing task", taskID)
	//get service by name
	services, err := c.docker.ListServices(docker.ListServicesOptions{})
	if err != nil {
		return errors.New("getting service list", err)
	}
	for _, svc := range services {
		//name is task id
		if strings.Contains(svc.Spec.Name, c.serviceName(taskID)) {
			if err := c.docker.RemoveService(docker.RemoveServiceOptions{ID: svc.ID}); err != nil {
				return errors.New("deleting service " + taskID, err)
			}
			return nil
		}
	}
	return errors.New("service with name " + taskID + " not found in service list " + fmt.Sprintf("%+v", services), nil)
}

func (c *Client) GetStatuses() ([]*mesosproto.TaskStatus, error) {
	statuses := []*mesosproto.TaskStatus{}
	logrus.Debug("getting status for all services")
	services, err := c.docker.ListServices(docker.ListServicesOptions{})
	if err != nil {
		return errors.New("getting service list", err)
	}
	tasks, err := c.docker.ListTasks(docker.ListTasksOptions{})
	if err != nil {
		return errors.New("getting task list", err)
	}
	for _, svc := range services {
		//service belongs to the rpi
		if strings.Contains(svc.Spec.Name, c.rpiName) {
			for _, task := range tasks {
				//find the only task that should exist for this service
				if task.ServiceID == svc.ID {
					statuses = append(statuses, c.convertToStatus(svc, task))
					break
				}
			}
		}
	}

	logrus.Debug("got status list", statuses)
	return statuses, nil
}

func (c *Client) serviceName(taskID string) string {
	return c.rpiName + "+" + taskID
}

func (c *Client) taskID(serviceName string) string {
	return strings.TrimPrefix(serviceName, c.rpiName + "+")
}

func (c *Client) convertToServiceSpec(task *lxtypes.Task, nodeName string) docker.CreateServiceOptions {
	logrus.Debug("converting task", task, "to service spec")

	nodeID := nodeName
	name := c.serviceName(task.TaskId)
	volumes := task.Container.GetVolumes()
	image := task.Container.Docker.GetImage()

	//env
	env := []string{}
	if task.Command.Environment != nil {
		for _, pair := range task.Command.Environment.Variables {
			env = append(env, pair.GetName() + "=" + pair.GetValue())
		}
	}

	//args or command for container
	var cmd []string
	args := []string{}
	if task.Command != nil && task.Command.Value != nil {
		if binary := *task.Command.Value; binary != "" {
			cmd = strings.Split(binary, " ")
		}
	}
	if task.Command.Arguments != nil {
		args = append(args, task.Command.Arguments...)
	}

	//mounts
	mounts := []mount.Mount{}
	for _, volume := range volumes {
		ro := false
		if volume.GetMode() == mesosproto.Volume_RO {
			ro = true
		}
		mount := mount.Mount{
			//todo: support persistent volumes
			Type: mount.TypeBind,
			Source: volume.GetHostPath(),
			Target: volume.GetContainerPath(),
			ReadOnly: ro,
		}
		mounts = append(mounts, mount)
	}

	//ports
	ports := []swarm.PortConfig{}
	for _, mapping := range task.Container.Docker.GetPortMappings() {
		port := swarm.PortConfig{
			Protocol: swarm.PortConfigProtocol(mapping.GetProtocol()), //tcp or udp
			TargetPort: mapping.GetContainerPort(),
			PublishedPort: mapping.GetHostPort(),
		}
		ports = append(ports, port)
	}

	//resources
	resources := &swarm.Resources{
		NanoCPUs: int64(task.Cpus * 1e9),
		MemoryBytes: int64(task.Mem << 20),
	}
	return docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: name,
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Image: image,
					Command: cmd,
					Args: args,
					Env: env,
					Mounts: mounts,
				},
				Resources: &swarm.ResourceRequirements{
					Limits: resources,
					Reservations: resources,
				},
				RestartPolicy: &swarm.RestartPolicy{
					Condition: swarm.RestartPolicyConditionNone,
				},
				Placement: &swarm.Placement{
					Constraints: []string{"node.id==" + nodeID},
				},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Ports: ports,
			},
		},
	}
}

func (c *Client) convertToStatus(service swarm.Service, task swarm.Task) mesosproto.StatusUpdate {
	taskID := c.taskID(service.Spec.Name)

	message := task.Status.Message
	var mesosState mesosproto.TaskState
	switch task.Status.State {
	case swarm.TaskStatePending:
		fallthrough
	case swarm.TaskStatePreparing:
		mesosState = mesosproto.TaskState_TASK_STARTING
	case swarm.TaskStateReady:
		fallthrough
	case swarm.TaskStateRunning:
		mesosState = mesosproto.TaskState_TASK_RUNNING
	case swarm.TaskStateComplete:
		mesosState = mesosproto.TaskState_TASK_FINISHED
	case swarm.TaskStateRejected:
		fallthrough
	case swarm.TaskStateFailed:
		mesosState = mesosproto.TaskState_TASK_FAILED
	default:
		mesosState = mesosproto.TaskState_TASK_ERROR
	}
	return &mesosproto.TaskStatus{
		TaskId:  &mesosproto.TaskID{Value: proto.String(taskID)},
		State:   &mesosState,
		Message: proto.String(message),
		SlaveId: &mesosproto.SlaveID{Value: proto.String(task.NodeID)},
	}
}