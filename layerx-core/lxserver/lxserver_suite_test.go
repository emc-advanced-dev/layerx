package lxserver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLxserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lxserver Suite")
}
