package layerx_tpi_api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLayerxTpiApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LayerxTpiApi Suite")
}
