package task_launcher_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTaskLauncher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TaskLauncher Suite")
}
