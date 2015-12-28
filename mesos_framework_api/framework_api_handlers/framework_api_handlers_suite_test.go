package framework_api_handlers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFrameworkApiHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FrameworkApiHandlers Suite")
}
