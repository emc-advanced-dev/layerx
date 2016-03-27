package main
import (
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-core_v2/lxstate"
	"flag"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-core_v2/lxserver"
	"fmt"
	"github.com/layer-x/layerx-core_v2/main_loop"
	"github.com/layer-x/layerx-core_v2/task_launcher"
	"github.com/go-martini/martini"
	"github.com/layer-x/layerx-core_v2/health_checker"
)

func purgeState() error {
	return lxdatabase.Rmdir("/state", true)
}

func main(){
	portPtr := flag.Int("port", 6666, "port to run core on")
	etcdUrlPtr := flag.String("etcd", "127.0.0.1:4001", "url of etcd cluster")
	purgePtr := flag.Bool("purge", false, "purge ETCD state")
	debugPtr := flag.Bool("debug", false, "Run Layer-X in debug mode")
	flag.Parse()

	if *debugPtr {
		lxlog.ActiveDebugMode()
	}

	lxlog.Infof(logrus.Fields{
		"port": *portPtr,
		"etcd": *etcdUrlPtr,
	}, "Booting Layer-X Core...")

	state := lxstate.NewState()
	err := state.InitializeState("http://"+*etcdUrlPtr)
	if err != nil {
		lxlog.Fatalf(logrus.Fields{
			"etcd": *etcdUrlPtr,
		}, "Failed to initialize Layer-X State")
	}
	if *purgePtr {
		err = purgeState()
		if err != nil {
			lxlog.Fatalf(logrus.Fields{
				"etcd": *etcdUrlPtr,
			}, "Failed to purge Layer-X State")
		}
		err = state.InitializeState("http://"+*etcdUrlPtr)
		if err != nil {
			lxlog.Fatalf(logrus.Fields{
				"etcd": *etcdUrlPtr,
			}, "Failed to initialize Layer-X State")
		}
	}

	lxlog.Infof(logrus.Fields{
		"port": *portPtr,
		"etcd": *etcdUrlPtr,
	}, "Layer-X Core Initialized. Waiting for registration of TPI and RPI...")

	driverErrc := make(chan error)
	mainServer := lxmartini.QuietMartini()
	coreServerWrapper := lxserver.NewLayerXCoreServerWrapper(state, mainServer, driverErrc)

	mainServer = coreServerWrapper.WrapServer()

	mainServer.Use(martini.Static("web"))

	go mainServer.RunOnAddr(fmt.Sprintf(":%v", *portPtr))

	clearRpisAndResources(state)

	taskLauncher := task_launcher.NewTaskLauncher(state)
	healthChecker := health_checker.NewHealthChecker(state)
	go main_loop.MainLoop(taskLauncher, healthChecker, state, driverErrc)
	lxlog.Infof(logrus.Fields{
	}, "Layer-X Server initialized successfully.")

	for {
		err = <- driverErrc
		if err != nil {
			lxlog.Errorf(logrus.Fields{"error": err},
				"Layer-X Core had an error!")
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
			lxlog.Warnf(logrus.Fields{"err": err, "node": node}, "retreiving resource pool for node")
			continue
		}
		resources, _ := nodeResourcePool.GetResources()
		for _, resource := range resources {
			nodeResourcePool.DeleteResource(resource.Id)
		}
	}
}