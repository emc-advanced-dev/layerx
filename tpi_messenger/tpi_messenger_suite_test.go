package tpi_messenger_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTpiMessenger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TpiMessenger Suite")
}
