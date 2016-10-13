/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/layerx/layerx-swarm-rpi/server"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/golang/protobuf/proto"
	"github.com/layer-x/layerx-commons/lxutils"
	"github.com/mesos/mesos-go/mesosproto"
	"net"
	"time"
	"github.com/emc-advanced-dev/layerx/layerx-swarm-rpi/swarm"
)

const (
	rpi_name = "Docker-Swarm-RPI-0.0.0"
)

var (
	port       = flag.String("port", "4000", "port to run on")
	name       = flag.String("name", rpi_name, "unique name to use for this rpi")
	layerX     = flag.String("layerx", "", "address:port for layerx core")
	localIpStr = flag.String("localip", "", "broadcast ip for the rpi")
	debug      = flag.Bool("debug", false, "verbose logging")

	statusUpdateInterval = time.Millisecond * 250 //1/4 second per status update
)

func main() {
	flag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	//register to layer x core
	localip := net.ParseIP(*localIpStr)
	if localip == nil {
		var err error
		localip, err = lxutils.GetLocalIp()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatalf("retrieving local ip")
		}
	}
	core := &layerx_rpi_client.LayerXRpi{
		CoreURL: *layerX,
		RpiName: *name,
	}
	logrus.WithFields(logrus.Fields{
		"rpi_url": fmt.Sprintf("%s:%v", localip.String(), *port),
	}).Infof("registering to layerx")

	if err := core.RegisterRpi(*name, fmt.Sprintf("%s:%v", localip.String(), *port)); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err.Error(),
			"layerx_url": *layerX,
		}).Fatal("registering to layerx")
	}

	client, err := swarm.NewClient(*name)
	if err != nil {
		logrus.Fatal("failed starting swarm client: ", err)
	}

	go func() {
		statusUpdatesForever(core, client)
	}()

	server.Start(*port, client, core)
}

func statusUpdatesForever(core *layerx_rpi_client.LayerXRpi, client *swarm.Client) {
	for {
		statuses, err := client.GetStatuses()
		if err != nil {
			logrus.Error("failed retrieving k8s task status updates", err)
			continue
		}
		resources, err := client.FetchResources()
		if err != nil {
			logrus.Error("failed retrieving k8s resource list", err)
			continue
		}
		killedStatuses, err := updatesForKilledTasks(core, statuses, resources)
		if err != nil {
			logrus.Error("failed to get killed tasks", err)
			continue
		}
		statuses = append(statuses, killedStatuses...)
		for _, status := range statuses {
			go func() {
				if err := core.SubmitStatusUpdate(status); err != nil {
					logrus.Error("failed submitting status updates to core", err)
				}
			}()
		}
		time.Sleep(statusUpdateInterval)
	}
}

func updatesForKilledTasks(core *layerx_rpi_client.LayerXRpi, notKilledStatuses []*mesosproto.TaskStatus, ownedResources []*lxtypes.Resource) ([]*mesosproto.TaskStatus, error) {
	killedStatuses := []*mesosproto.TaskStatus{}
	nodes, err := core.GetNodes()
	if err != nil {
		return nil, errors.New("failed checking the lx list of nodes", err)
	}
	for _, node := range nodes {
		ownsNode := false
		for _, resource := range ownedResources {
			if resource.NodeId == node.Id {
				ownsNode = true
				break
			}
		}
		if !ownsNode {
			continue
		}
		for _, task := range node.GetTasks() {
			taskKilled := true
			for _, status := range notKilledStatuses {
				if status.GetTaskId().GetValue() == task.TaskId {
					taskKilled = false
					break
				}
			}
			if taskKilled {
				killedStatuses = append(killedStatuses, newPodKilledStatus(task))
			}
		}
	}
	return killedStatuses, nil
}

func newPodKilledStatus(task *lxtypes.Task) *mesosproto.TaskStatus {
	var mesosState = mesosproto.TaskState_TASK_KILLED
	return &mesosproto.TaskStatus{
		TaskId:  &mesosproto.TaskID{Value: proto.String(task.TaskId)},
		State:   &mesosState,
		Message: proto.String("Task Killed or Lost"),
		SlaveId: &mesosproto.SlaveID{Value: proto.String(task.NodeId)},
	}
}
