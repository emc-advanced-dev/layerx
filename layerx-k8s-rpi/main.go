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
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"flag"
	"k8s.io/client-go/1.4/tools/clientcmd"
	"github.com/Sirupsen/logrus"
	"k8s.io/client-go/1.4/kubernetes"
	"fmt"
	"net"
	"github.com/layer-x/layerx-commons/lxutils"
	"github.com/emc-advanced-dev/layerx/layerx-k8s-rpi/kube"
	"github.com/emc-advanced-dev/layerx/layerx-k8s-rpi/server"
)

const (
	rpi_name = "Kubernetes-RPI-0.0.0"
)

var (
	kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
	port = flag.String("port", "4000", "port to run on")
	layerX = flag.String("layerx", "", "address:port for layerx core")
	localIpStr = flag.String("localip", "", "broadcast ip for the rpi")
	debug = flag.Bool("debug", false, "verbose logging")

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
	rpiClient := &layerx_rpi_client.LayerXRpi{
		CoreURL: *layerX,
		RpiName: rpi_name,
	}
	logrus.WithFields(logrus.Fields{
		"rpi_url": fmt.Sprintf("%s:%v", localip.String(), *port),
	}).Infof("registering to layerx")

	if err := rpiClient.RegisterRpi(rpi_name, fmt.Sprintf("%s:%v", localip.String(), *port)); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err.Error(),
			"layerx_url": *layerX,
		}).Errorf("registering to layerx")
	}



	//initialize kube client
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		logrus.Fatal(err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatal(err)
	}

	kubeClient := kube.NewClient(clientset)

	server.Start(*port, kubeClient)

	//for {
	//	pods, err := clientset.Core().Pods("").List(api.ListOptions{})
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	//	time.Sleep(10 * time.Second)
	//}
}
