package main
import (
	"flag"
	"github.com/Sirupsen/logrus"
"github.com/layer-x/layerx-commons/lxlog"
"github.com/layer-x/layerx-commons/lxutils"
	"fmt"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/layer-x/layerx-mesos-tpi_v2/driver"
	"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api"
	"github.com/layer-x/layerx-core_v2/layerx_tpi"
)

func main () {
	port := flag.Int("port", 3030, "listening port for mesos tpi, default: 3000")
	debug := flag.String("debug", "false", "turn on debugging, default: false")
	layerX := flag.String("layerx", "", "layer-x url, e.g. \"10.141.141.10:3000\"")

	flag.Parse()

	if *debug == "true" {
		lxlog.ActiveDebugMode()
		lxlog.Debugf(logrus.Fields{}, "debugging activated")
	}

	if *layerX == "" {
		lxlog.Fatalf(logrus.Fields{
		}, "-layerx flag not set")
	}

	localip, err := lxutils.GetLocalIp()
	if err != nil {
		lxlog.Fatalf(logrus.Fields{
			"error": err.Error(),
		}, "retrieving local ip")
	}
	masterUpidString := fmt.Sprintf("master@%s:%v", localip.String(), *port)
	masterUpid, err := mesos_data.UPIDFromString(masterUpidString)
	if err != nil {
		lxlog.Fatalf(logrus.Fields{
			"error": err.Error(),
			"masterUpidString": masterUpidString,
		}, "generating master upid")
	}

	actionQueue := lxactionqueue.NewActionQueue()
	driver := driver.NewMesosTpiDriver(actionQueue)
	frameworkManager := framework_manager.NewFrameworkManager(masterUpid)
	tpi := &layerx_tpi.LayerXTpi{
		CoreURL: *layerX,
	}

	masterServer := mesos_master_api.NewMesosApiServer(tpi, actionQueue, frameworkManager)
	errc := make(chan error)
	go masterServer.RunMasterServer(*port, masterUpidString, errc)
	go driver.Run()

	lxlog.Infof(logrus.Fields{
		"port":          *port,
		"layer-x-url":   *layerX,
		"upid":            masterUpidString,
		"driver":    		driver,
	}, "Layerx Mesos TPI initialized...")

	err = <- errc
	if err != nil {
		lxlog.Fatalf(logrus.Fields{
			"error": err.Error(),
		}, "Mesos server failed")
	}

}