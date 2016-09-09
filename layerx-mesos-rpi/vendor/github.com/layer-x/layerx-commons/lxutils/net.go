package lxutils

import (
	"github.com/layer-x/layerx-commons/lxerrors"
	"net"
	"strconv"
	"strings"
)

func GetLocalIp() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	// handle err
	for _, i := range ifaces {
		if i.Name == "eth1" {
			addrs, err := i.Addrs()
			if err != nil {
				return nil, err
			}
			// handle err
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					return v.IP, nil
				case *net.IPAddr:
					return v.IP, nil
				}
				// process IP address
			}
		}
	}
	return nil, lxerrors.New("Could not find IP in network interfaces", nil)
}

//Taken from https://groups.google.com/forum/#!topic/golang-nuts/v4eJ5HK3stI
// Convert net.IP to uint32
func IpToI(ipnr net.IP) uint32 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}
