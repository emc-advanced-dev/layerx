package layerx_rpi_api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLayerxRpiApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LayerxRpiApi Suite")
}
