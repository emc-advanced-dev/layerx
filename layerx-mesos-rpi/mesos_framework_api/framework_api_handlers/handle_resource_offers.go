package framework_api_handlers
import (
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func HandleResourceOffers(lxRpi *layerx_rpi_client.LayerXRpi, offers []*mesosproto.Offer) error {
	for _, offer := range offers {
		resource := lxtypes.NewResourceFromMesos(offer)
		err := lxRpi.SubmitResource(resource)
		if err != nil {
			return lxerrors.New("failed to submit resource "+resource.Id+" to layerx core", err)
		}
	}
	return nil
}