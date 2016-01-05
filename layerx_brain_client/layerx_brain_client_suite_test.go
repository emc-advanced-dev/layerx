package layerx_brain_client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLayerxBrainClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LayerxBrainClient Suite")
}
