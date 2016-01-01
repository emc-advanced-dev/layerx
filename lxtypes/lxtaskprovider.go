package lxtypes

type TaskProvider struct {
	Id string `json:"id"`
	//Source should be some kind of contact info to
	//reach a Task Provider. e.g. in the case of Mesos,
	//this should contain UPID of a framework so
	//the TPI can send messages to the framework
	Source string `json:"source"`
	//indicates whether task provider
	//has a failover timeout
	FailoverTimeout float64 `json:"failover_timeout"`
}
