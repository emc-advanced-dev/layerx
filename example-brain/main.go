package main

import (
	"github.com/emc-advanced-dev/pkg/errors"
	"time"
	"github.com/Sirupsen/logrus"
	"flag"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_brain_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
)

var core *layerx_brain_client.LayerXBrainClient

func main() {
	layerXAddress := flag.String("core", "", "ip:port of layerx core")
	flag.Parse()

	//main loop of scheduling / reading the cluster state
	core = &layerx_brain_client.LayerXBrainClient{CoreURL: *layerXAddress}
	if err := mainLoop(); err != nil {
		logrus.Fatal("scheduling loop failed", err)
	}
}

func mainLoop() error {
	for {
		if err := scheduleOnce(); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
}

func scheduleOnce() error {
	logrus.Info("requesting nodes & tasks from ", core.CoreURL)
	nodes, err := core.GetNodes()
	if err != nil {
		return errors.New("retrieving node list", err)
	}
	tasks, err := core.GetPendingTasks()
	if err != nil {
		return errors.New("retrieving task list", err)
	}
	logrus.Info("current state of the cluster: ", nodes)
	if len(tasks) < 1 {
		logrus.Info("no nodes available for scheduling.")
		return nil
	}
	if len(tasks) < 1 {
		logrus.Info("nothing pending to schedule.")
		return nil
	}
	for _, task := range tasks {
		node, err := pickNode(nodes, task)
		if err != nil {
			logrus.WithFields(logrus.Fields{"task": task, "nodes": nodes}).Error(err)
			continue
		}
		if err := core.AssignTasks(node.Id, task.TaskId); err != nil {
			return errors.New("failed to assign task "+task.TaskId+" to node "+node.Id, err)
		}
		//change the local representation of the node for purposes of pickNode()
		node.AddTask(task)

		logrus.Info("assigned task ",task, " to node ", node, "!")
	}
	return nil
}

func pickNode(nodes []*lxtypes.Node, task *lxtypes.Task) (*lxtypes.Node, error) {
	//find the node with the fewest running tasks
	//ignore nodes that don't have available capacity
	var selectedNode *lxtypes.Node
	for _, node := range nodes {
		if node.GetFreeCpus() < task.Cpus || node.GetFreeDisk() < task.Disk || node.GetFreeMem() < task.Mem {
			//dismiss this node
			continue
		}
		//choose this node if:
		// we haven't selected a node yet (first node in loop)
		// the node we're inspecting has fewer running tasks than the node we selected
		if selectedNode == nil || len(selectedNode.GetTasks()) > len(node.GetTasks()) {
			selectedNode = node
		}
	}
	if selectedNode == nil {
		return nil, errors.New("failed to find a suitable node for task "+task.TaskId, nil)
	}
	return selectedNode, nil
}
