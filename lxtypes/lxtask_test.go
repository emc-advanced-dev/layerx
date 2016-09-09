package lxtypes_test

import (
	. "github.com/emc-advanced-dev/layerx-core/lxtypes"

	"github.com/emc-advanced-dev/layerx-core/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lxtask", func() {
	Describe("NewTaskFromMesos()", func() {
		It("converts back from a mesos task", func() {
			mesosTask := fakes.FakeMesosTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			task := NewTaskFromMesos(mesosTask)
			Expect(task.Command).To(Equal(mesosTask.Command))
			Expect(task.Container).To(Equal(mesosTask.Container))
			Expect(task.Name).To(Equal(mesosTask.GetName()))
			Expect(task.TaskId).To(Equal(mesosTask.GetTaskId().GetValue()))
			Expect(task.NodeId).To(Equal(mesosTask.GetSlaveId().GetValue()))
		})
	})
	Describe("ToMesos()", func() {
		It("converts to a mesos task", func() {
			mesosTask := fakes.FakeMesosTask("fake_task_id", "fake_task_name", "fake_slave_id", "echo FAKE_COMMAND")
			//need this to add the extra labels
			task := NewTaskFromMesos(mesosTask)
			mesosTask = task.ToMesos()
			//
			task = NewTaskFromMesos(mesosTask)
			convertedTask := task.ToMesos()
			Expect(convertedTask).To(Equal(mesosTask))
		})
	})
})
