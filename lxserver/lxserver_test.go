package lxserver_test

import (
	. "github.com/layer-x/layerx-core_v2/lxserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-core_v2/layerx_rpi_client"
	"github.com/layer-x/layerx-core_v2/layerx_tpi_client"
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"github.com/layer-x/layerx-commons/lxmartini"
	"github.com/layer-x/layerx-core_v2/driver"
	"github.com/layer-x/layerx-core_v2/lxstate"
"github.com/layer-x/layerx-commons/lxlog"
	"fmt"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"github.com/layer-x/layerx-core_v2/fakes"
)


func PurgeState() {
	lxdatabase.Rmdir("/state", true)
}

var _ = Describe("Lxserver", func() {
	var lxRpiClient *layerx_rpi_client.LayerXRpi
	var lxTpiClient *layerx_tpi_client.LayerXTpi
	var state *lxstate.State

	Describe("setup", func(){
		It("sets up for the tests", func(){
			lxRpiClient = &layerx_rpi_client.LayerXRpi{
				CoreURL: "127.0.0.1:6677",
			}
			lxTpiClient = &layerx_tpi_client.LayerXTpi{
				CoreURL: "127.0.0.1:6677",
			}

			actionQueue := lxactionqueue.NewActionQueue()
			state = lxstate.NewState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			coreServerWrapper := NewLayerXCoreServerWrapper(state, actionQueue)
			driver := driver.NewLayerXDriver(actionQueue)

			m := coreServerWrapper.WrapServer(lxmartini.QuietMartini(), make(chan error))
			go m.RunOnAddr(fmt.Sprintf(":6677"))
			go fakes.RunFakeTpiServer("127.0.0.1:6677", 6688, make(chan error))
			go driver.Run()
			lxlog.ActiveDebugMode()
		})
	})

	Describe("RegisterTpi", func() {
		It("adds the Tpi URL to the LX state", func() {
			PurgeState()
			err := lxTpiClient.RegisterTpi("127.0.0.1:6688")
			Expect(err).To(BeNil())
			tpiUrl, err := state.GetTpi()
			Expect(err).To(BeNil())
			Expect(tpiUrl).To(Equal("127.0.0.1:6688"))
		})
	})

	Describe("RegisterRpi", func() {
		It("adds the Rpi URL to the LX state", func() {
			PurgeState()
			err := lxRpiClient.RegisterRpi("127.0.0.1:6699")
			Expect(err).To(BeNil())
			rpiUrl, err := state.GetRpi()
			Expect(err).To(BeNil())
			Expect(rpiUrl).To(Equal("127.0.0.1:6699"))
		})
	})

	Describe("RegisterTaskProvider", func() {
		It("adds the task provider to the LX state", func() {
			PurgeState()
			fakeTaskProvider := fakes.FakeTaskProvider("fake_framework", "ff@fakeip:fakeport")
			err := lxTpiClient.RegisterTaskProvider(fakeTaskProvider)
			Expect(err).To(BeNil())
			taskProvider, err := state.TaskProviderPool.GetTaskProvider("fake_framework")
			Expect(err).To(BeNil())
			Expect(taskProvider).To(Equal(fakeTaskProvider))
		})
	})
})
