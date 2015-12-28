package lxtypes_test

import (
	. "github.com/layer-x/layerx-core_v2/lxtypes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-core_v2/fakes"
)

var _ = Describe("Lxnode", func() {
	Describe("AddResource", func() {
		Context("nodeId of resource matches nodeId of node", func() {
			It("adds the resource to the node", func() {
				fakeOffer := fakes.FakeOffer("fake_offer_id", "fake_slave_id")
				fakeResource := NewResourceFromMesos(fakeOffer)
				fakeNode := NewNode(fakeResource.NodeId)
				err := fakeNode.AddResource(fakeResource)
				Expect(err).To(BeNil())
				Expect(fakeNode.GetResource("fake_offer_id")).To(Equal(fakeResource))
			})
		})
		Context("nodeId of resource does not match nodeId of node", func() {
			It("returns an error", func() {
				fakeOffer := fakes.FakeOffer("fake_offer_id", "fake_slave_id")
				fakeResource := NewResourceFromMesos(fakeOffer)
				fakeNode := NewNode("other_" + fakeResource.NodeId)
				err := fakeNode.AddResource(fakeResource)
				Expect(err).NotTo(BeNil())
				Expect(fakeNode.GetResource("fake_offer_id")).To(BeNil())
			})
		})
		Context("resourceId is already found on node", func() {
			It("returns an error", func() {
				fakeOffer := fakes.FakeOffer("fake_offer_id", "fake_slave_id")
				fakeResource := NewResourceFromMesos(fakeOffer)
				fakeNode := NewNode(fakeResource.NodeId)
				err := fakeNode.AddResource(fakeResource)
				Expect(err).To(BeNil())
				err = fakeNode.AddResource(fakeResource)
				Expect(err).NotTo(BeNil())
			})
		})
	})
	Describe("GetResource", func() {
		Context("resource is found on node", func() {
			It("returns the resource for that resource", func() {
				fakeOffer := fakes.FakeOffer("fake_offer_id", "fake_slave_id")
				fakeResource := NewResourceFromMesos(fakeOffer)
				fakeNode := NewNode(fakeResource.NodeId)
				err := fakeNode.AddResource(fakeResource)
				Expect(err).To(BeNil())
				Expect(fakeNode.GetResource("fake_offer_id")).To(Equal(fakeResource))
				task := fakeNode.GetResource(fakeResource.Id)
				Expect(task).To(Equal(fakeResource))
			})
		})
		Context("resource is not found on node", func() {
			It("returns nil", func() {
				fakeNode := NewNode("fake_node_id")
				fakeOffer := fakes.FakeOffer("fake_offer_id", "fake_slave_id")
				fakeResource := NewResourceFromMesos(fakeOffer)
				task := fakeNode.GetResource(fakeResource.Id)
				Expect(task).To(BeNil())
			})
		})
	})
	Describe("FlushResources", func() {
		It("removes the all resources from the node", func() {
			fakeOffer := fakes.FakeOffer("fake_offer_id", "fake_slave_id")
			fakeResource := NewResourceFromMesos(fakeOffer)
			fakeNode := NewNode(fakeResource.NodeId)
			err := fakeNode.AddResource(fakeResource)
			Expect(err).To(BeNil())
			Expect(fakeNode.GetResource("fake_offer_id")).To(Equal(fakeResource))
			fakeNode.FlushResources()
			Expect(fakeNode.GetResource("fake_offer_id")).To(BeNil())
		})
	})
	Describe("AddTask", func() {
		Context("taskid is not found on node", func() {
			It("adds the task to the node", func() {
				fakeNode := NewNode("fake_node_id")
				fakeMesosTask := fakes.FakeMesosTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask := NewTaskFromMesos(fakeMesosTask)
				err := fakeNode.AddTask(fakeTask)
				Expect(err).To(BeNil())
				task := fakeNode.GetTask(fakeTask.TaskId)
				Expect(task).To(Equal(fakeTask))
			})
		})
		Context("taskId is already found on node", func() {
			It("returns an error", func() {
				fakeNode := NewNode("fake_node_id")
				fakeMesosTask := fakes.FakeMesosTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask := NewTaskFromMesos(fakeMesosTask)
				err := fakeNode.AddTask(fakeTask)
				Expect(err).To(BeNil())
				task := fakeNode.GetTask(fakeTask.TaskId)
				Expect(task).To(Equal(fakeTask))
				err = fakeNode.AddTask(fakeTask)
				Expect(err).NotTo(BeNil())
			})
		})
	})
	Describe("ModifyTask", func() {
		Context("taskid is found on node", func() {
			It("modifies the existing task on the node", func() {
				fakeNode := NewNode("fake_node_id")
				fakeMesosTask := fakes.FakeMesosTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask := NewTaskFromMesos(fakeMesosTask)
				err := fakeNode.AddTask(fakeTask)
				Expect(err).To(BeNil())
				task := fakeNode.GetTask(fakeTask.TaskId)
				Expect(task).To(Equal(fakeTask))
				fakeModifiedMesosTask := fakes.FakeMesosTask("fake_task_id", "other_fake_task_name", "other_fake_slave_id", "echo other_FAKE_COMMAND")
				fakeModifiedTask := NewTaskFromMesos(fakeModifiedMesosTask)
				err = fakeNode.ModifyTask(fakeModifiedTask)
				Expect(err).To(BeNil())
				task = fakeNode.GetTask(fakeTask.TaskId)
				Expect(task).NotTo(Equal(fakeTask))
				Expect(task).To(Equal(fakeModifiedTask))
			})
		})
		Context("taskId is not found on node", func() {
			It("returns an error", func() {
				fakeNode := NewNode("fake_node_id")
				fakeMesosTask := fakes.FakeMesosTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask := NewTaskFromMesos(fakeMesosTask)
				err := fakeNode.ModifyTask(fakeTask)
				Expect(err).NotTo(BeNil())
				task := fakeNode.GetTask(fakeTask.TaskId)
				Expect(task).To(BeNil())
			})
		})
	})
	Describe("GetTask", func() {
		Context("taskid is found on node", func() {
			It("returns the Task for that taskid", func() {
				fakeNode := NewNode("fake_node_id")
				fakeMesosTask := fakes.FakeMesosTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask := NewTaskFromMesos(fakeMesosTask)
				err := fakeNode.AddTask(fakeTask)
				Expect(err).To(BeNil())
				task := fakeNode.GetTask(fakeTask.TaskId)
				Expect(task).To(Equal(fakeTask))
			})
		})
		Context("taskId is not found on node", func() {
			It("returns nil", func() {
				fakeNode := NewNode("fake_node_id")
				fakeMesosTask := fakes.FakeMesosTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask := NewTaskFromMesos(fakeMesosTask)
				task := fakeNode.GetTask(fakeTask.TaskId)
				Expect(task).To(BeNil())
			})
		})
	})
	Describe("RemoveTask", func() {
		Context("taskid is found on node", func() {
			It("removes the task from the node", func() {
				fakeNode := NewNode("fake_node_id")
				fakeMesosTask := fakes.FakeMesosTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask := NewTaskFromMesos(fakeMesosTask)
				err := fakeNode.AddTask(fakeTask)
				Expect(err).To(BeNil())
				task := fakeNode.GetTask(fakeTask.TaskId)
				Expect(task).To(Equal(fakeTask))
				err = fakeNode.RemoveTask(fakeTask.TaskId)
				Expect(err).To(BeNil())
				task = fakeNode.GetTask(fakeTask.TaskId)
				Expect(task).To(BeNil())
			})
		})
		Context("taskId is not found on node", func() {
			It("returns error", func() {
				fakeNode := NewNode("fake_node_id")
				fakeMesosTask := fakes.FakeMesosTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
				fakeTask := NewTaskFromMesos(fakeMesosTask)
				err := fakeNode.RemoveTask(fakeTask.TaskId)
				Expect(err).NotTo(BeNil())
			})
		})
	})

})
