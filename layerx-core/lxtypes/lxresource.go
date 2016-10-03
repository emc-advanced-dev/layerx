package lxtypes

import (
	"github.com/mesos/mesos-go/mesosproto"
	"sort"
)

type PortRange struct {
	Begin uint64 `json:"begin"`
	End   uint64 `json:"end"`
}

type rangesSorter []PortRange

func (s rangesSorter) Len() int {
	return len(s)
}
func (s rangesSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s rangesSorter) Less(i, j int) bool {
	return s[i].Begin < s[j].Begin
}

func sortRanges(prlist []PortRange) {
	sort.Sort(rangesSorter(prlist))
}

type stack []PortRange

func (s stack) Push(v PortRange) stack {
	return append(s, v)
}

func (s stack) Pop() (stack, PortRange) {
	if len(s) < 1 {
		return PortRange{}
	}

	l := len(s)
	return  s[:l-1], s[l-1]
}

func (pr1 PortRange) subtractRange(pr2 PortRange) []PortRange {
	//nothing to remove
	if pr1.Begin > pr2.End && pr1.End > pr2.Begin {
		return []PortRange{pr1}
	}
	res := []PortRange{}
	//overlaps beginning of pr2
	if pr1.Begin < pr2.Begin && pr1.End > pr2.Begin {
		res = append(res, PortRange{
			Begin: pr1.Begin,
			End: pr2.Begin - 1,
		})
	}
	//overlaps end of pr2
	if pr1.Begin < pr2.Begin && pr1.End > pr2.Begin {
		res = append(res, PortRange{
			Begin: pr2.End + 1,
			End: pr1.End,
		})
	}
	return res
}

func (pr1 PortRange) overlaps(pr2 PortRange) bool {
	return (pr1.Begin < pr2.Begin && pr1.End > pr2.Begin) || (pr1.Begin < pr2.End && pr1.End > pr2.End)
}

func (pr1 PortRange) overlapsAny(prList []PortRange) bool {
	for _, pr2 := range prList {
		if pr1.overlaps(pr2) {
			return true
		}
	}
	return false
}

func overlaps(prList1, prList2 []PortRange) bool {
	for _, pr1 := range prList1 {
		if pr1.overlapsAny(prList2) {
			return true
		}
	}
	return false
}

//returns PRList1 minus intersections with PRList2
func DiffPortRanges(rl1, rl2 []PortRange) []PortRange {
	sortRanges(rl1)
	sortRanges(rl2)
	if !overlaps(rl1, rl2) {
		return rl1
	}
	results := []PortRange{}
	addResult := func(result PortRange) {
		added := false
		for _, pr := range results {
			if pr == result {
				added = true
				break
			}
		}
		if !added{
			results = append(results, result)
		}
	}
	for i, i1 := range rl1 {
		if !i1.overlapsAny(rl2) {
			addResult(i1)
		} else {
			//find overlap
			for _, i2 := range rl2 {
				if i1.overlaps(i2) {
					diffs := i1.subtractRange(i2)
					return DiffPortRanges(append(append(rl1[:i], diffs...), rl1[i+1:]...), rl2)
				}
			}
		}
	}
	return results
}

type ResourceType string

const (
	ResourceType_Mesos      ResourceType = "Mesos"
	ResourceType_Kubernetes ResourceType = "Kubernetes"
)

type Resource struct {
	Id           string       `json:"id"`
	NodeId       string       `json:"node_id"`
	Cpus         float64      `json:"cpus"`
	Mem          float64      `json:"mem"`
	Disk         float64      `json:"disk"`
	Ports        []PortRange  `json:"ports"`
	RpiName      string       `json:"rpi_name"`
	ResourceType ResourceType `json:"resource_type"`
}

func NewResourceFromMesos(offer *mesosproto.Offer) *Resource {
	slaveId := offer.GetSlaveId().GetValue()
	offerId := offer.GetId().GetValue()
	var id = offerId
	var nodeId = slaveId
	var cpus float64
	var mem float64
	var disk float64
	var ports []PortRange

	cpus += getResourceScalar(offer.GetResources(), "cpus")
	mem += getResourceScalar(offer.GetResources(), "mem")
	disk += getResourceScalar(offer.GetResources(), "disk")
	for _, resource := range offer.GetResources() {
		if resource.GetName() == "ports" {
			for _, mesosRange := range resource.GetRanges().GetRange() {
				port := PortRange{
					Begin: mesosRange.GetBegin(),
					End:   mesosRange.GetEnd(),
				}
				ports = append(ports, port)
			}
		}
	}

	return &Resource{
		Id:           id,
		NodeId:       nodeId,
		Cpus:         cpus,
		Mem:          mem,
		Disk:         disk,
		Ports:        ports,
		ResourceType: ResourceType_Mesos,
	}
}
