package lx_core_helpers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLxCoreHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LxCoreHelpers Suite")
}
