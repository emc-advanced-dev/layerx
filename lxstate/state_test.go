package lxstate_test

import (
	. "github.com/layer-x/layerx-core_v2/lxstate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-commons/lxdatabase"
)

func PurgeState() {
	lxdatabase.Rmdir("/state", true)
}

var _ = Describe("State", func() {
	Describe("InitializeState(etcdUrl)", func() {
		It("initializes client (lxdb), creates folders for nodes, tasks, statuses, tps", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			rootContents, err := lxdatabase.GetSubdirectories("/state")
			Expect(err).To(BeNil())
			Expect(rootContents).To(ContainElement("/state/nodes"))
			Expect(rootContents).To(ContainElement("/state/pending_tasks"))
			Expect(rootContents).To(ContainElement("/state/staging_tasks"))
			Expect(rootContents).To(ContainElement("/state/task_providers"))
			Expect(rootContents).To(ContainElement("/state/statuses"))
		})
	})
	Describe("Set/GetTpiUrl(tpiUrl)", func() {
		It("sets and gets the tpiurl", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			err = state.SetTpi("fake_url")
			Expect(err).To(BeNil())
			tpiUrl, err := state.GetTpi()
			Expect(err).To(BeNil())
			Expect(tpiUrl).To(Equal("fake_url"))
		})
	})
	Describe("Set/GetRpiUrl(tpiUrl)", func() {
		It("sets and gets the rpiurl", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			err = state.SetRpi("fake_url")
			Expect(err).To(BeNil())
			rpiUrl, err := state.GetRpi()
			Expect(err).To(BeNil())
			Expect(rpiUrl).To(Equal("fake_url"))
		})
	})
})
