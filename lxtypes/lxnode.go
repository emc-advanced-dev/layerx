package lxtypes
import "github.com/layer-x/layerx-commons/lxerrors"

type Node interface {
	AddResource(resource *Resource) error
	GetResource(resourceId string) *Resource
	GetResources() []*Resource
	FlushResources()
	AddTask(task *Task) error
	GetTask(taskId string) *Task
	ModifyTask(task *Task) error
	RemoveTask(taskId string) error
	GetTotalCpus() float64
	GetTotalMem() float64
	GetTotalPorts() []PortRange
	GetFreeCpus() float64
	GetFreeMem() float64
	GetFreeDisk() float64
	GetFreePorts() []PortRange
}

type node struct {
	Id string        `json:"id"`
	Resources map[string]*Resource `json:"resources"`
	RunningTasks map[string]*Task `json:"tasks"`
}

func NewNode(nodeId string) *node {
	return &node{
		Id: nodeId,
		Resources: make(map[string]*Resource),
		RunningTasks: make(map[string]*Task),
	}
}

func (n *node) AddResource(resource *Resource) error {
	if _, ok := n.Resources[resource.Id]; ok {
		return lxerrors.New("resource " + resource.Id + " already found on node " + n.Id, nil)
	}
	if n.Id != resource.NodeId {
		return lxerrors.New("attempted to add resource " + resource.Id + " with node id " + resource.NodeId + " to node " + n.Id, nil)
	}
	n.Resources[resource.Id] = resource
	return nil
}

func (n *node) GetResource(resourceId string) *Resource {
	return n.Resources[resourceId]
}

func (n *node) GetResources() []*Resource {
	resources := []*Resource{}
	for _, resource := range n.Resources {
		resources = append(resources, resource)
	}
	return resources
}

func (n *node) FlushResources() {
	n.Resources = make(map[string]*Resource)
}

func (n *node) AddTask(task *Task) error {
	if _, ok := n.RunningTasks[task.TaskId]; ok {
		return lxerrors.New("task " + task.TaskId + " already found on node " + n.Id, nil)
	}
	n.RunningTasks[task.TaskId] = task
	return nil
}

func (n *node) GetTask(taskId string) *Task {
	return n.RunningTasks[taskId]
}

func (n *node) ModifyTask(task *Task) error {
	if _, ok := n.RunningTasks[task.TaskId]; !ok {
		return lxerrors.New("task " + task.TaskId + " not found on node " + n.Id, nil)
	}
	n.RunningTasks[task.TaskId] = task
	return nil
}

func (n *node) RemoveTask(taskId string) error {
	if _, ok := n.RunningTasks[taskId]; !ok {
		return lxerrors.New("task " + taskId + " not found on node " + n.Id, nil)
	}
	delete(n.RunningTasks, taskId)
	return nil
}

func (n *node) GetTotalCpus() float64 {
	cpus := n.GetFreeCpus()
	for _, task := range n.RunningTasks {
		cpus += task.Cpus
	}
	return cpus
}

func (n *node) GetTotalMem() float64 {
	mem := n.GetFreeMem()
	for _, task := range n.RunningTasks {
		mem += task.Mem
	}
	return mem
}

func (n *node) GetTotalDisk() float64 {
	disk := n.GetFreeDisk()
	for _, task := range n.RunningTasks {
		disk += task.Disk
	}
	return disk
}

func (n *node) GetTotalPorts() []PortRange {
	totalPorts := n.GetFreePorts()
	for _, task := range n.RunningTasks {
		for _, port := range task.Ports {
			portRange := PortRange{
				Begin: port.Begin,
				End:   port.End,
			}
			totalPorts = append(totalPorts, portRange)
		}
	}
	return totalPorts
}

func (n *node) GetFreeCpus() float64 {
	var cpus float64
	for _, resource := range n.Resources {
		cpus += resource.Cpus
	}
	return cpus}

func (n *node) GetFreeMem() float64 {
	var mem float64
	for _, resource := range n.Resources {
		mem += resource.Mem
	}
	return mem
}

func (n *node) GetFreeDisk() float64 {
	var disk float64
	for _, resource := range n.Resources {
		disk += resource.Disk
	}
	return disk
}

func (n *node) GetFreePorts() []PortRange {
	freePorts := []PortRange{}
	for _, resource := range n.Resources {
		for _, port := range resource.Ports {
			portRange := PortRange{
				Begin: port.Begin,
				End:   port.End,
			}
			freePorts = append(freePorts, portRange)
		}
	}
	return freePorts
}