package state_test


import (
	. "github.com/layer-x/layerx-core_v2/state"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-commons/lxdatabase"
	"encoding/json"
	"github.com/layer-x/layerx-core_v2/fakes"
	"github.com/layer-x/layerx-core_v2/lxtypes"
)

var _ = Describe("ResourcePool", func() {
	Describe("GetResource(resourceId)", func(){
		It("returns the resource if it exists, else returns err", func(){
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			fakeResourcePool := TempResourcePoolFunction("/state/fake_node_pool/fake_node_id_1/fake_resource_pool", "fake_node_id_1")
			fakeResource := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_1", "fake_node_id_1"))
			resource, err := fakeResourcePool.GetResource(fakeResource.Id)
			Expect(err).NotTo(BeNil())
			Expect(resource).To(BeNil())
			err = fakeResourcePool.AddResource(fakeResource)
			Expect(err).To(BeNil())
			resource, err = fakeResourcePool.GetResource(fakeResource.Id)
			Expect(err).To(BeNil())
			Expect(resource).To(Equal(fakeResource))
		})
	})
	Describe("AddResource", func(){
		Context("the resource is new", func(){
			It("adds the resource to etcd state", func(){
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeResourcePool := TempResourcePoolFunction("/state/fake_node_pool/fake_node_id_1/fake_resource_pool", "fake_node_id_1")
				fakeResource := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_1", "fake_node_id_1"))
				err = fakeResourcePool.AddResource(fakeResource)
				Expect(err).To(BeNil())
				expectedResourceJsonBytes, err := json.Marshal(fakeResource)
				Expect(err).To(BeNil())
				expectedResourceJson := string(expectedResourceJsonBytes)
				actualResourceJson, err := lxdatabase.Get(fakeResourcePool.GetKey() + "/"+fakeResource.Id)
				Expect(err).To(BeNil())
				Expect(actualResourceJson).To(Equal(expectedResourceJson))
			})
		})
		Context("the resource is not new", func(){
			It("returns an error", func(){
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeResourcePool := TempResourcePoolFunction("/state/fake_node_pool/fake_node_id_1/fake_resource_pool", "fake_node_id_1")
				fakeResource := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_1", "fake_node_id_1"))
				err = fakeResourcePool.AddResource(fakeResource)
				Expect(err).To(BeNil())
				err = fakeResourcePool.AddResource(fakeResource)
				Expect(err).NotTo(BeNil())
			})
		})
		Context("the resource does not belong to the node", func(){
			It("returns an error", func(){
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeResourcePool := TempResourcePoolFunction("/state/fake_node_pool/fake_node_id_1/fake_resource_pool", "fake_node_id_1")
				fakeResource := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_1", "fake_node_id_2"))
				err = fakeResourcePool.AddResource(fakeResource)
				Expect(err).NotTo(BeNil())
			})
		})
	})
	Describe("ModifyResource", func(){
		Context("the exists", func(){
			It("modifies the resource", func(){
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeResourcePool := TempResourcePoolFunction("/state/fake_node_pool/fake_node_id_1/fake_resource_pool", "fake_node_id_1")
				fakeResource := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_1", "fake_node_id_1"))
				err = fakeResourcePool.AddResource(fakeResource)
				Expect(err).To(BeNil())
				fakeResource.Mem = 666
				fakeResource.Cpus = 666
				fakeResource.Disk = 666
				err = fakeResourcePool.ModifyResource(fakeResource.Id, fakeResource)
				Expect(err).To(BeNil())
				expectedResourceJsonBytes, err := json.Marshal(fakeResource)
				Expect(err).To(BeNil())
				expectedResourceJson := string(expectedResourceJsonBytes)
				actualResourceJson, err := lxdatabase.Get(fakeResourcePool.GetKey() + "/"+fakeResource.Id)
				Expect(err).To(BeNil())
				Expect(actualResourceJson).To(Equal(expectedResourceJson))
			})
		})
		Context("the resource doest exist", func(){
			It("returns an error", func(){
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeResourcePool := TempResourcePoolFunction("/state/fake_node_pool/fake_node_id_1/fake_resource_pool", "fake_node_id_1")
				fakeResource := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_1", "fake_node_id_1"))
				err = fakeResourcePool.ModifyResource(fakeResource.Id, fakeResource)
				Expect(err).NotTo(BeNil())
			})
		})
	})
	Describe("GetResources()", func(){
		It("returns all known resources in the pool", func(){
			state := NewState()
			state.InitializeState("http://127.0.0.1:4001")
			PurgeState()
			err := state.InitializeState("http://127.0.0.1:4001")
			Expect(err).To(BeNil())
			fakeResourcePool := TempResourcePoolFunction("/state/fake_node_pool/fake_node_id_1/fake_resource_pool", "fake_node_id_1")
			fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_1", "fake_node_id_1"))
			fakeResource2 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_2", "fake_node_id_1"))
			fakeResource3 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_3", "fake_node_id_1"))
			err = fakeResourcePool.AddResource(fakeResource1)
			Expect(err).To(BeNil())
			err = fakeResourcePool.AddResource(fakeResource2)
			Expect(err).To(BeNil())
			err = fakeResourcePool.AddResource(fakeResource3)
			Expect(err).To(BeNil())
			resources, err := fakeResourcePool.GetResources()
			Expect(err).To(BeNil())
			Expect(resources[fakeResource1.Id]).To(Equal(fakeResource1))
			Expect(resources[fakeResource2.Id]).To(Equal(fakeResource2))
			Expect(resources[fakeResource3.Id]).To(Equal(fakeResource3))
		})
	})
	Describe("DeleteResource(resourceId)", func(){
		Context("resource exists", func(){
			It("deletes the resource", func(){
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeResourcePool := TempResourcePoolFunction("/state/fake_node_pool/fake_node_id_1/fake_resource_pool", "fake_node_id_1")
				fakeResource1 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_1", "fake_node_id_1"))
				fakeResource2 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_2", "fake_node_id_1"))
				fakeResource3 := lxtypes.NewResourceFromMesos(fakes.FakeOffer("fake_resource_id_3", "fake_node_id_1"))
				err = fakeResourcePool.AddResource(fakeResource1)
				Expect(err).To(BeNil())
				err = fakeResourcePool.AddResource(fakeResource2)
				Expect(err).To(BeNil())
				err = fakeResourcePool.AddResource(fakeResource3)
				Expect(err).To(BeNil())
				err = fakeResourcePool.DeleteResource(fakeResource1.Id)
				Expect(err).To(BeNil())
				resources, err := fakeResourcePool.GetResources()
				Expect(err).To(BeNil())
				Expect(resources[fakeResource1.Id]).To(BeNil())
				Expect(resources[fakeResource2.Id]).To(Equal(fakeResource2))
				Expect(resources[fakeResource3.Id]).To(Equal(fakeResource3))
			})
		})
		Context("resource does not exist", func(){
			It("throws error", func(){
				state := NewState()
				state.InitializeState("http://127.0.0.1:4001")
				PurgeState()
				err := state.InitializeState("http://127.0.0.1:4001")
				Expect(err).To(BeNil())
				fakeResourcePool := TempResourcePoolFunction("/state/fake_node_pool/fake_node_id_1/fake_resource_pool", "fake_node_id_1")
				err = fakeResourcePool.DeleteResource("nonexistent_resource_id")
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
