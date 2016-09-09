package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_tpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/framework_manager"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/layerx_tpi_api"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/mesos_master_api"
	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-commons/lxutils"
)

func main() {
	port := flag.Int("port", 3000, "listening port for mesos tpi")
	debug := flag.String("debug", "false", "turn on debugging, default: false")
	layerX := flag.String("layerx", "", "layer-x url, e.g. \"10.141.141.10:5000\"")
	localIpStr := flag.String("localip", "", "binding address for the rpi")

	flag.Parse()

	if *debug == "true" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugf("debugging activated")
	}

	if *layerX == "" {
		logrus.WithFields(logrus.Fields{}).Fatalf("-layerx flag not set")
	}

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
	masterUpidString := fmt.Sprintf("master@%s:%v", localip.String(), *port)
	masterUpid, err := mesos_data.UPIDFromString(masterUpidString)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":            err.Error(),
			"masterUpidString": masterUpidString,
		}).Fatalf("generating master upid")
	}

	frameworkManager := framework_manager.NewFrameworkManager(masterUpid)
	tpi := &layerx_tpi_client.LayerXTpi{
		CoreURL: *layerX,
	}

	masterServerWrapper := mesos_master_api.NewMesosApiServerWrapper(tpi, frameworkManager)
	tpiServerWrapper := layerx_tpi_api.NewTpiApiServerWrapper(tpi, frameworkManager)
	errc := make(chan error)
	tpiServer := lxmartini.QuietMartini()
	tpiServer = masterServerWrapper.WrapWithMesos(tpiServer, masterUpidString, errc)
	tpiServer = tpiServerWrapper.WrapWithTpi(tpiServer, masterUpidString, errc)

	go tpiServer.RunOnAddr(fmt.Sprintf(":%v", *port))

	err = tpi.RegisterTpi(fmt.Sprintf("%s:%v", localip, *port))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err.Error(),
			"layerx_url": *layerX,
		}).Errorf("registering to layerx")
	}

	logrus.WithFields(logrus.Fields{
		"port":        *port,
		"layer-x-url": *layerX,
		"upid":        masterUpidString,
	}).Infof("Layerx Mesos TPI initialized...")

	for {
		err = <-errc
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("Mesos tpi experienced a failure!")
		}
	}

}
