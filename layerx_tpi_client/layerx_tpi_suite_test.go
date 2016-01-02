package layerx_tpi_client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLayerxTpi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LayerxTpi Suite")
}
