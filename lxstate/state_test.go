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
})
