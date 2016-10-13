package swarm

import "time"

type swarmNode struct {
	ID string `json:"ID"`
	Version struct {
		   Index int `json:"Index"`
	   } `json:"Version"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Spec struct {
		   Role string `json:"Role"`
		   Availability string `json:"Availability"`
	   } `json:"Spec"`
	Description struct {
		   Hostname string `json:"Hostname"`
		   Platform struct {
				    Architecture string `json:"Architecture"`
				    OS string `json:"OS"`
			    } `json:"Platform"`
		   Resources struct {
				    NanoCPUs int `json:"NanoCPUs"`
				    MemoryBytes int `json:"MemoryBytes"`
			    } `json:"Resources"`
		   Engine struct {
				    EngineVersion string `json:"EngineVersion"`
				    Labels struct {
							  Provider string `json:"provider"`
						  } `json:"Labels"`
				    Plugins []struct {
					    Type string `json:"Type"`
					    Name string `json:"Name"`
				    } `json:"Plugins"`
			    } `json:"Engine"`
	   } `json:"Description"`
	Status struct {
		   State string `json:"State"`
	   } `json:"Status"`
}