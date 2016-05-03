package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-commons/lxutils"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-mesos-rpi_v2/layerx_rpi_api"
	"github.com/layer-x/layerx-mesos-rpi_v2/mesos_framework_api"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
)

const rpi_name = "Mesos-RPI-0.0.0"

func main() {
	port := flag.Int("port", 4040, "listening port for mesos rpi, default: 2999")
	master := flag.String("master", "127.0.0.1:5050", "url of mesos master")
	debug := flag.String("debug", "false", "turn on debugging, default: false")
	layerX := flag.String("layerx", "", "layer-x url, e.g. \"10.141.141.10:3000\"")
	localIpStr := flag.String("localip", "", "binding address for the rpi")
	rpiName := flag.String("name", rpi_name, "name to use to register to layerx")
	flag.Parse()

	if *debug == "true" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugf("debugging activated")
	}

	localip := net.ParseIP(*localIpStr)
	if localip == nil {
		var err error
		localip, err = lxutils.GetLocalIp()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatalf( "retrieving local ip")
		}
	}

	rpiFramework := prepareFrameworkInfo(*layerX)
	rpiClient := &layerx_rpi_client.LayerXRpi{
		CoreURL: *layerX,
		RpiName: *rpiName,
	}

	logrus.WithFields(logrus.Fields{
		"rpi_url": fmt.Sprintf("%s:%v", localip.String(), *port),
	}).Infof("registering to layerx")

	err := rpiClient.RegisterRpi(*rpiName, fmt.Sprintf("%s:%v", localip.String(), *port))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err.Error(),
			"layerx_url": *layerX,
		}).Errorf( "registering to layerx")
	}

	rpiScheduler := mesos_framework_api.NewRpiMesosScheduler(rpiClient)

	config := scheduler.DriverConfig{
		Scheduler:  rpiScheduler,
		Framework:  rpiFramework,
		Master:     *master,
		Credential: (*mesosproto.Credential)(nil),
	}

	go func() {
		driver, err := scheduler.NewMesosSchedulerDriver(config)
		if err != nil {
			err = lxerrors.New("initializing mesos schedulerdriver", err)
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"mesos_url": *master,
			}).Fatalf( "error initializing mesos schedulerdriver")
		}
		status, err := driver.Run()
		if err != nil {
			err = lxerrors.New("Framework stopped with status "+status.String(), err)
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"mesos_url": *master,
			}).Fatalf( "error running mesos schedulerdriver")
			return
		}
	}()
	mesosSchedulerDriver := rpiScheduler.GetDriver()
	rpiServerWrapper := layerx_rpi_api.NewRpiApiServerWrapper(rpiClient, mesosSchedulerDriver)
	errc := make(chan error)
	m := rpiServerWrapper.WrapWithRpi(lxmartini.QuietMartini(), errc)
	go m.RunOnAddr(fmt.Sprintf(":%v", *port))

	logrus.WithFields(logrus.Fields{
		"config": config,
	}).Infof("Layer-X Mesos RPI Initialized...")

	for {
		err = <-errc
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err}).Errorf( "LayerX Mesos RPI Failed!")
		}
	}
}

func prepareFrameworkInfo(layerxUrl string) *mesosproto.FrameworkInfo {
	return &mesosproto.FrameworkInfo{
		User: proto.String(""),
		//		Id: &mesosproto.FrameworkID{
		//			Value: proto.String("lx_mesos_rpi_framework_3"),
		//		},
		FailoverTimeout: proto.Float64(0),
		Name:            proto.String("Layer-X Mesos RPI Framework"),
		WebuiUrl:        proto.String(layerxUrl),
	}
}
