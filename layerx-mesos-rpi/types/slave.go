package types

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"strings"
	"strconv"
	"github.com/emc-advanced-dev/pkg/errors"
)

type Slaves struct {
	Slaves []Slave `json:"slaves"`
}

type Slave struct {
	ID                    string `json:"id"`
	Pid                   string `json:"pid"`
	Hostname              string `json:"hostname"`
	RegisteredTime        float64 `json:"registered_time"`
	ReregisteredTime      float64 `json:"reregistered_time"`
	Resources             struct {
				      Cpus  float64 `json:"cpus"`
				      Disk  float64 `json:"disk"`
				      Mem   float64 `json:"mem"`
				      Ports string `json:"ports"`
			      } `json:"resources"`
	UsedResources         struct {
				      Cpus  float64 `json:"cpus"`
				      Disk  float64 `json:"disk"`
				      Mem   float64 `json:"mem"`
				      Ports string `json:"ports"`
			      } `json:"used_resources"`
	OfferedResources      struct {
				      Cpus float64 `json:"cpus"`
				      Disk float64 `json:"disk"`
				      Mem  float64 `json:"mem"`
			      } `json:"offered_resources"`
	ReservedResources     struct {
			      } `json:"reserved_resources"`
	UnreservedResources   struct {
				      Cpus  float64 `json:"cpus"`
				      Disk  float64 `json:"disk"`
				      Mem   float64 `json:"mem"`
				      Ports string `json:"ports"`
			      } `json:"unreserved_resources"`
	Attributes            struct {
			      } `json:"attributes"`
	Active                bool `json:"active"`
	Version               string `json:"version"`
	ReservedResourcesFull struct {
			      } `json:"reserved_resources_full"`
	UsedResourcesFull     []struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Scalar struct {
			       Value float64 `json:"value"`
		       } `json:"scalar,omitempty"`
		Role   string `json:"role"`
		Ranges struct {
			       Range []struct {
				       Begin int `json:"begin"`
				       End   int `json:"end"`
			       } `json:"range"`
		       } `json:"ranges,omitempty"`
	} `json:"used_resources_full"`
	OfferedResourcesFull  []interface{} `json:"offered_resources_full"`
}

func (s *Slave) ToResource() (*lxtypes.Resource, error) {
	totalPorts, err := parsePortsString(s.Resources.Ports)
	if err != nil {
		return nil, errors.New("failed to parse ports string "+s.Resources.Ports, err)
	}
	usedPorts, err := parsePortsString(s.UsedResources.Ports)
	if err != nil {
		return nil, errors.New("failed to parse ports string "+s.UsedResources.Ports, err)
	}
	return &lxtypes.Resource{
		Id:     s.ID,
		NodeId: s.ID,
		Cpus:   s.Resources.Cpus - s.UsedResources.Cpus,
		Mem:    s.Resources.Mem - s.UsedResources.Mem,
		Disk:   s.Resources.Disk - s.UsedResources.Disk,
		Ports:  lxtypes.DiffPortRanges(totalPorts, usedPorts),
		ResourceType: lxtypes.ResourceType_Mesos,
	}, nil
}

func parsePortsString(portString string) ([]lxtypes.PortRange, error) {
	ports := []lxtypes.PortRange{}
	portString = strings.TrimPrefix(portString, "[")
	portString = strings.TrimSuffix(portString, "]")
	if len(portString) < 1 {
		return ports, nil
	}
	portRanges := strings.Split(portString, ", ")
	for _, portRange := range portRanges {
		split := strings.Split(portRange, "-")
		begin, err := strconv.Atoi(split[0])
		if err != nil {
			return nil, errors.New("could not parse "+split[0]+" as integer", err)
		}
		end, err := strconv.Atoi(split[1])
		if err != nil {
			return nil, errors.New("could not parse "+split[1]+" as integer", err)
		}
		ports = append(ports, lxtypes.PortRange{
			Begin: uint64(begin),
			End: uint64(end),
		})
	}
	return ports, nil
}