package main
import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/layer-x/layerx-commons/lxutils"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
	"github.com/gogo/protobuf/proto"
	"github.com/layer-x/layerx-mesos-rpi_v2/mesos_framework_api"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/layer-x/layerx-commons/lxerrors"
	"github.com/layer-x/layerx-mesos-rpi_v2/layerx_rpi_api"
	"github.com/layer-x/layerx-mesos-rpi_v2/driver"
	"github.com/layer-x/layerx-commons/lxmartini"
	"fmt"
)

func main() {
	port := flag.Int("port", 4040, "listening port for mesos rpi, default: 2999")
	master := flag.String("master", "127.0.0.1:5050", "url of mesos master")
	debug := flag.String("debug", "false", "turn on debugging, default: false")
	layerX := flag.String("layerx", "", "layer-x url, e.g. \"10.141.141.10:3000\"")
	flag.Parse()

	if *debug == "true" {
		lxlog.ActiveDebugMode()
		lxlog.Debugf(logrus.Fields{}, "debugging activated")
	}

	localip, err := lxutils.GetLocalIp()
	if err != nil {
		lxlog.Fatalf(logrus.Fields{
			"error": err.Error(),
		}, "retrieving local ip")
	}
	actionQueue := lxactionqueue.NewActionQueue()

	rpiFramework := prepareFrameworkInfo(*layerX)
	rpiClient := &layerx_rpi_client.LayerXRpi{
		CoreURL: *layerX,
	}

	lxlog.Infof(logrus.Fields{
		"rpi_url": fmt.Sprintf("%s:%v", localip.String(), *port),
	}, "registering to layerx")

	err = rpiClient.RegisterRpi(fmt.Sprintf("%s:%v", localip.String(), *port))
	if err != nil {
		lxlog.Fatalf(logrus.Fields{
			"error": err.Error(),
			"layerx_url": *layerX,
		}, "registering to layerx")
	}

	rpiScheduler := mesos_framework_api.NewRpiMesosScheduler(rpiClient, actionQueue)

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
			lxlog.Fatalf(logrus.Fields{
				"error":     err,
				"mesos_url": *master,
			}, "error initializing mesos schedulerdriver")
		}
		status, err := driver.Run()
		if err != nil {
			err = lxerrors.New("Framework stopped with status " + status.String(), err)
			lxlog.Fatalf(logrus.Fields{
				"error":     err,
				"mesos_url": *master,
			}, "error running mesos schedulerdriver")
			return
		}
	}()
	mesosSchedulerDriver := rpiScheduler.GetDriver()
	rpiServerWrapper := layerx_rpi_api.NewRpiApiServerWrapper(rpiClient, mesosSchedulerDriver, actionQueue)
	driver := driver.NewMesosRpiDriver(actionQueue)
	errc := make(chan error)
	m := rpiServerWrapper.WrapWithRpi(lxmartini.QuietMartini(), errc)
	go m.RunOnAddr(fmt.Sprintf(":%v", *port))
	go driver.Run()

	lxlog.Infof(logrus.Fields{
		"config": config,
	}, "Layer-X Mesos RPI Initialized...")

	err = <-errc
	if err != nil {
		lxlog.Fatalf(logrus.Fields{"error": err}, "LayerX Mesos RPI Failed!")
	}
}

func prepareFrameworkInfo(layerxUrl string) *mesosproto.FrameworkInfo {
	return &mesosproto.FrameworkInfo{
		User: proto.String(""),
		Id: &mesosproto.FrameworkID{
			Value: proto.String("lx_mesos_rpi_framework"),
		},
		FailoverTimeout: proto.Float64(15),
		Name: proto.String("Layer-X Mesos RPI Framework"),
		WebuiUrl:        proto.String(layerxUrl),
	}
}