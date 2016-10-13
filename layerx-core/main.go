package main

import (
	"flag"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/bindata"
	"github.com/emc-advanced-dev/layerx/layerx-core/health_checker"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxserver"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/layerx/layerx-core/main_loop"
	"github.com/emc-advanced-dev/layerx/layerx-core/task_launcher"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/emc-advanced-dev/pkg/logger"
)

func purgeState() error {
	return lxdatabase.Rmdir("/state", true)
}

func main() {
	portPtr := flag.Int("port", 5000, "port to run core on")
	etcdUrlPtr := flag.String("etcd", "127.0.0.1:4001", "url of etcd cluster")
	purgePtr := flag.Bool("purge", false, "purge ETCD state on boot")
	debugPtr := flag.Bool("debug", false, "verbose logging")
	flag.Parse()

	if *debugPtr {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.AddHook(logger.LoggerNameHook{"CORE"})

	logrus.WithFields(logrus.Fields{
		"port": *portPtr,
		"etcd": *etcdUrlPtr,
	}).Infof("Booting Layer-X Core...")


	state := lxstate.NewState()
	err := state.InitializeState("http://" + *etcdUrlPtr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"etcd": *etcdUrlPtr,
		}).Fatalf("Failed to initialize Layer-X State")
	}
	if *purgePtr {
		err = purgeState()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"etcd": *etcdUrlPtr,
			}).Fatalf("Failed to purge Layer-X State")
		}
		err = state.InitializeState("http://" + *etcdUrlPtr)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"etcd": *etcdUrlPtr,
			}).Fatalf("Failed to initialize Layer-X State")
		}
	}

	logrus.WithFields(logrus.Fields{
		"port": *portPtr,
		"etcd": *etcdUrlPtr,
	}).Infof("Layer-X Core Initialized. Waiting for registration of TPI and RPI...")

	driverErrc := make(chan error)
	mainServer := lxmartini.QuietMartini()
	coreServerWrapper := lxserver.NewLayerXCoreServerWrapper(state, mainServer, driverErrc)

	mainServer = coreServerWrapper.WrapServer()

	mainServer.Use(bindata.Static())

	go mainServer.RunOnAddr(fmt.Sprintf(":%v", *portPtr))

	clearRpisAndResources(state)

	taskLauncher := task_launcher.NewTaskLauncher(state)
	healthChecker := health_checker.NewHealthChecker(state)
	go main_loop.MainLoop(taskLauncher, healthChecker, state, driverErrc)
	logrus.WithFields(logrus.Fields{}).Infof("Layer-X Server initialized successfully.")

	for {
		err = <-driverErrc
		if err != nil {
			logrus.WithError(err).Errorf("Layer-X Core had an error!")
		}
	}
}

func clearRpisAndResources(state *lxstate.State) {
	//clear previous rpis
	oldRpis, _ := state.RpiPool.GetRpis()
	for _, rpi := range oldRpis {
		state.RpiPool.DeleteRpi(rpi.Name)
	}
	//clear previous resources
	oldNodes, _ := state.NodePool.GetNodes()
	for _, node := range oldNodes {
		nodeResourcePool, err := state.NodePool.GetNodeResourcePool(node.Id)
		if err != nil {
			logrus.WithFields(logrus.Fields{"err": err, "node": node}).Warnf("retreiving resource pool for node")
			continue
		}
		resources, _ := nodeResourcePool.GetResources()
		for _, resource := range resources {
			nodeResourcePool.DeleteResource(resource.Id)
		}
	}
}
