package framework_manager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFrameworkManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FrameworkManager Suite")
}
