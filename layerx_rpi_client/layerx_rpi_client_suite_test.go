package layerx_rpi_client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLayerxRpiClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LayerxRpiClient Suite")
}
