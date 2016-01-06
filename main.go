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
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-mesos-tpi_v2/layerx_tpi_api"
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
	tpi := &layerx_tpi_client.LayerXTpi{
		CoreURL: *layerX,
	}

	masterServerWrapper := mesos_master_api.NewMesosApiServerWrapper(tpi, actionQueue, frameworkManager)
	tpiServerWrapper := layerx_tpi_api.NewTpiApiServerWrapper(tpi, actionQueue, frameworkManager)
	errc := make(chan error)
	tpiServer := lxmartini.QuietMartini()
	tpiServer = masterServerWrapper.WrapWithMesos(tpiServer, masterUpidString, errc)
	tpiServer = tpiServerWrapper.WrapWithTpi(tpiServer, masterUpidString, errc)

	go tpiServer.RunOnAddr(fmt.Sprintf(":%v",*port))
	go driver.Run()

	err = tpi.RegisterTpi(fmt.Sprintf("%s:%v",localip, *port))
	if err != nil {
		lxlog.Fatalf(logrus.Fields{
			"error": err.Error(),
			"layerx_url": *layerX,
		}, "registering to layerx")
	}

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