package lxstate_test

import (
	. "github.com/emc-advanced-dev/layerx/layerx-core/lxstate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"encoding/json"
	"github.com/layer-x/layerx-commons/lxdatabase"
)

func fakeLXRpi(name, url string) *layerx_rpi_client.RpiInfo {
	return &layerx_rpi_client.RpiInfo{
		Name: name,
		Url: url,
	}
}

var _ = Describe("RpiPool", func() {
	Describe("GetRpi(rpiName)", func() {
		It("returns the rpi if it exists, else returns err", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			fakeRpi := fakeLXRpi("fake_rpi_id_1", "http://fake.url")
			rpi, err := state.RpiPool.GetRpi(fakeRpi.Name)
			Expect(err).NotTo(BeNil())
			Expect(rpi).To(BeNil())
			err = state.RpiPool.AddRpi(fakeRpi)
			Expect(err).To(BeNil())
			rpi, err = state.RpiPool.GetRpi(fakeRpi.Name)
			Expect(err).To(BeNil())
			Expect(rpi).To(Equal(fakeRpi))
		})
	})
	Describe("AddRpi", func() {
		Context("the rpi is new", func() {
			It("adds the rpi to etcd state", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeRpi := fakeLXRpi("fake_rpi_id_1", "http://fake.url")

				err = state.RpiPool.AddRpi(fakeRpi)
				Expect(err).To(BeNil())
				expectedRpiJsonBytes, err := json.Marshal(fakeRpi)
				Expect(err).To(BeNil())
				expectedRpiJson := string(expectedRpiJsonBytes)
				actualRpiJson, err := lxdatabase.Get(state.RpiPool.GetKey() + "/" + fakeRpi.Name)
				Expect(err).To(BeNil())
				Expect(actualRpiJson).To(Equal(expectedRpiJson))
			})
		})
		Context("the rpi is not new", func() {
			It("replaces the old entry with a new one", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeRpi := fakeLXRpi("fake_rpi_id_1", "http://fake.url")

				err = state.RpiPool.AddRpi(fakeRpi)
				Expect(err).To(BeNil())
				err = state.RpiPool.AddRpi(fakeRpi)
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("GetRpis()", func() {
		It("returns all known rpis in the pool", func() {
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			fakeRpi1 := fakeLXRpi("fake_rpi_id_1", "http://fake.url.1")
			fakeRpi2 := fakeLXRpi("fake_rpi_id_2", "http://fake.url.2")
			err = state.RpiPool.AddRpi(fakeRpi1)
			Expect(err).To(BeNil())
			err = state.RpiPool.AddRpi(fakeRpi2)
			Expect(err).To(BeNil())
			rpis, err := state.RpiPool.GetRpis()
			Expect(err).To(BeNil())
			Expect(rpis[fakeRpi1.Name]).To(Equal(fakeRpi1))
			Expect(rpis[fakeRpi2.Name]).To(Equal(fakeRpi2))
		})
	})
	Describe("DeleteRpi(rpiId)", func() {
		Context("rpi exists", func() {
			It("deletes the rpi", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeRpi1 := fakeLXRpi("fake_rpi_id_1", "http://fake.url")
				fakeRpi2 := fakeLXRpi("fake_rpi_id_2", "http://fake.url")
				fakeRpi3 := fakeLXRpi("fake_rpi_id_3", "http://fake.url")
				err = state.RpiPool.AddRpi(fakeRpi1)
				Expect(err).To(BeNil())
				err = state.RpiPool.AddRpi(fakeRpi2)
				Expect(err).To(BeNil())
				err = state.RpiPool.AddRpi(fakeRpi3)
				Expect(err).To(BeNil())
				err = state.RpiPool.DeleteRpi(fakeRpi1.Name)
				Expect(err).To(BeNil())
				rpis, err := state.RpiPool.GetRpis()
				Expect(err).To(BeNil())
				Expect(rpis[fakeRpi1.Name]).To(BeNil())
				Expect(rpis[fakeRpi2.Name]).To(Equal(fakeRpi2))
				Expect(rpis[fakeRpi3.Name]).To(Equal(fakeRpi3))
			})
		})
		Context("rpi does not exist", func() {
			It("throws error", func() {
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				err = state.RpiPool.DeleteRpi("nonexistent_rpi_id")
				Expect(err).NotTo(BeNil())
			})
		})
	})
})

