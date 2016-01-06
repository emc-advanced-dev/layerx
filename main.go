package main
import (
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-core_v2/lxstate"
	"flag"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"time"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-core_v2/lxserver"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/layer-x/layerx-core_v2/driver"
	"fmt"
	"github.com/layer-x/layerx-core_v2/main_loop"
	"github.com/layer-x/layerx-core_v2/task_launcher"
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

	actionQueue := lxactionqueue.NewActionQueue()
	driver := driver.NewLayerXDriver(actionQueue)
	//run driver
	go driver.Run()

	driverErrc := make(chan error)
	mainServer := lxmartini.QuietMartini()
	coreServerWrapper := lxserver.NewLayerXCoreServerWrapper(state, actionQueue, mainServer, "", "", driverErrc)

	mainServer = coreServerWrapper.WrapServer()
	go mainServer.RunOnAddr(fmt.Sprintf(":%v", *portPtr))


	rpiUrl := ""
	tpiUrl := ""
	for {
		tpiUrl, _ = state.GetTpi()
		rpiUrl, _ = state.GetRpi()
		if tpiUrl != "" && rpiUrl != "" {
			lxlog.Infof(logrus.Fields{
				"tpiUrl": tpiUrl,
				"rpiUrl": rpiUrl,
			}, "TPI and RPI have registered. Initializing Layer-X Server...")
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	taskLauncher := task_launcher.NewTaskLauncher(rpiUrl, state)
	go main_loop.MainLoop(actionQueue, taskLauncher, state, tpiUrl, rpiUrl, driverErrc)
	lxlog.Infof(logrus.Fields{
		"tpiUrl": tpiUrl,
		"rpiUrl": rpiUrl,
	}, "Layer-X Server initialized successfully.")

	err = <- driverErrc
	if err != nil {
		lxlog.Fatalf(logrus.Fields{"error": err},
			"Layer-X Core failed!")
	}
}