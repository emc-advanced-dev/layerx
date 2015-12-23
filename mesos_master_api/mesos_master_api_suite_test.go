package mesos_master_api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMesosMasterApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MesosMasterApi Suite")
}
