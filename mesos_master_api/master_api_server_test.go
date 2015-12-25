package mesos_master_api_test

import (
	. "github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
	"encoding/json"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/layer-x/layerx-mesos-tpi_v2/driver"
	"github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"
	"github.com/layer-x/layerx-mesos-tpi_v2/fakes"
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/layer-x/layerx-core_v2/layerx_tpi"
	core_fakes "github.com/layer-x/layerx-core_v2/fakes"
)

var _ = Describe("MasterApiServer", func() {
	actionQueue := lxactionqueue.NewActionQueue()
	fakeMasterUpid, _ := mesos_data.UPIDFromString("master@127.0.0.1:3031")
	frameworkManager := framework_manager.NewFrameworkManager(fakeMasterUpid)
	fakeTpi := &layerx_tpi.LayerXTpi{
		CoreURL: "127.0.0.1:34443",
	}
	masterServer := NewMesosApiServer(fakeTpi, actionQueue, frameworkManager)
	driver := driver.NewMesosTpiDriver(actionQueue)

	go masterServer.RunMasterServer(3031, "master@127.0.0.1:3031", make(chan error))
	go driver.Run()
	go fakes.RunFakeFrameworkServer("fakeframework", 3001)
	go core_fakes.RunFakeLayerXServer(nil, 34443)
	lxlog.ActiveDebugMode()

	Describe("GET " + GET_MASTER_STATE, func() {
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
	Describe("GET " + GET_MASTER_STATE_DEPRECATED, func() {
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
	Describe("POST {subscribe_call} " + MESOS_SCHEDULER_CALL, func() {
		It("registers the framework to layer-x core, returns \"FrameworkRegisteredMessage\" to framework", func() {
			fakeSubscribe := fakes.FakeSubscribeCall()
			headers := map[string]string{
				"Libprocess-From": "fakeframework@127.0.0.1:3001",
			}
			resp, _, err := lxhttpclient.Post("127.0.0.1:3031", MESOS_SCHEDULER_CALL, headers, fakeSubscribe)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(202))
		})
	})
})
