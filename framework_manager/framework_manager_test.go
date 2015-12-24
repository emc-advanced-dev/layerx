package framework_manager_test

import (
	. "github.com/layer-x/layerx-mesos-tpi_v2/framework_manager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
"github.com/layer-x/layerx-commons/lxlog"
"github.com/layer-x/layerx-mesos-tpi_v2/fakes"
	"github.com/layer-x/layerx-mesos-tpi_v2/mesos_master_api/mesos_data"
)

var _ = Describe("FrameworkManager", func() {
	go fakes.RunFakeFramework("fakeframework", 3001)
	lxlog.ActiveDebugMode()
	time.Sleep(3 * time.Second)


	Describe("Notify Framework is registered", func() {
		It("succesfully gets 202 response", func() {
			fakeMasterUpid, err := mesos_data.UPIDFromString("fakemesos@127.0.0.1:3031")
			Expect(err).To(BeNil())
			frameworkManager := NewFrameworkManager(fakeMasterUpid)
			frameworkUpid, err := mesos_data.UPIDFromString("fakeframework@127.0.0.1:3001")
			Expect(err).To(BeNil())
			err = frameworkManager.NotifyFrameworkRegistered("fakeframework",
				"fake_framework_id",
				frameworkUpid)
			Expect(err).To(BeNil())
		})
	})
})
