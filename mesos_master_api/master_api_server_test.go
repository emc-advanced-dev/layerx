package mesos_master_api_test

import (
	. "github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"encoding/json"
)

var _ = Describe("MasterApiServer", func() {

	go RunMasterServer(3031, "master@127.0.0.1:3031", make(chan error))

	Describe("GET "+GET_MASTER_STATE, func() {
		It("returns state of the faux master", func() {
			resp, data, err := lxhttpclient.Get("127.0.0.1:3031", GET_MASTER_STATE, nil)
			Expect(err).To(BeNil())
			var state mesos_data.MesosState
			Expect(resp.StatusCode).To(Equal(200))
			err = json.Unmarshal(data, &state)
			Expect(err).To(BeNil())
			Expect(state.Version).To(Equal("0.25.0"))
			Expect(state.Leader).To(Equal("master@127.0.0.1:3031"))
		})
	})
	Describe("GET "+GET_MASTER_STATE_DEPRECATED, func() {
		It("returns state of the faux master", func() {
			resp, data, err := lxhttpclient.Get("127.0.0.1:3031", GET_MASTER_STATE_DEPRECATED, nil)
			Expect(err).To(BeNil())
			var state mesos_data.MesosState
			Expect(resp.StatusCode).To(Equal(200))
			err = json.Unmarshal(data, &state)
			Expect(err).To(BeNil())
			Expect(state.Version).To(Equal("0.25.0"))
			Expect(state.Leader).To(Equal("master@127.0.0.1:3031"))
		})
	})
	Describe("POST {subscribe_call} "+MESOS_SCHEDULER_CALL, func() {
		It("Queues ", func() {
			resp, data, err := lxhttpclient.Post("127.0.0.1:3031", MESOS_SCHEDULER_CALL, nil, nil)
			Expect(err).To(BeNil())
			var state mesos_data.MesosState
			Expect(resp.StatusCode).To(Equal(200))
			err = json.Unmarshal(data, &state)
			Expect(err).To(BeNil())
			Expect(state.Version).To(Equal("0.25.0"))
			Expect(state.Leader).To(Equal("master@127.0.0.1:3031"))
		})
	})
})
