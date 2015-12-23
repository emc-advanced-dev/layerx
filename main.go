package main
import (
	"flag"
	"github.com/Sirupsen/logrus"
"github.com/layer-x/layerx-commons/lxlog"
"github.com/layer-x/layerx-commons/lxutils"
	"fmt"
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
	masterUpid := fmt.Sprintf("master@%s:%v", localip.String(), *port)

	lxlog.Infof(logrus.Fields{
		"port":          *port,
		"layer-x-url":   *layerX,
		"upid":    		masterUpid,
	}, "Layerx Mesos TPI initialized...")


}