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
	"github.com/emc-advanced-dev/layerx/layerx-k8s-rpi/kube"
	"github.com/emc-advanced-dev/layerx/layerx-k8s-rpi/server"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/golang/protobuf/proto"
	"github.com/layer-x/layerx-commons/lxutils"
	"github.com/mesos/mesos-go/mesosproto"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/tools/clientcmd"
	"net"
	"time"
	"strings"
)

const (
	rpi_name = "Kubernetes-RPI-0.0.0"
)

var (
	kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
	port       = flag.String("port", "4000", "port to run on")
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
		RpiName: rpi_name,
	}
	logrus.WithFields(logrus.Fields{
		"rpi_url": fmt.Sprintf("%s:%v", localip.String(), *port),
	}).Infof("registering to layerx")

	if err := core.RegisterRpi(rpi_name, fmt.Sprintf("%s:%v", localip.String(), *port)); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err.Error(),
			"layerx_url": *layerX,
		}).Fatal("registering to layerx")
	}

	//initialize kube client
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		logrus.Fatal("building config", err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatal("creating clientset", err)
	}

	kubeClient := kube.NewClient(clientset)
	if err := kubeClient.Init(); err != nil && !strings.Contains(err.Error(), "already exists"){
		logrus.Fatal("failed to initialize kubernetes namespace", err)
	}

	go func() {
		statusUpdatesForever(core, kubeClient)
	}()

	server.Start(*port, kubeClient, core)
}

func statusUpdatesForever(core *layerx_rpi_client.LayerXRpi, kubeClient *kube.Client) {
	for {
		statuses, err := kubeClient.GetStatuses()
		if err != nil {
			logrus.Error("failed retrieving k8s task status updates", err)
			continue
		}
		resources, err := kubeClient.FetchResources()
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
