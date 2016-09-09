package main_loop_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMainLoop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MainLoop Suite")
}
