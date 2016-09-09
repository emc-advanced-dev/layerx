package rpi_messenger_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRpiMessenger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RpiMessenger Suite")
}
