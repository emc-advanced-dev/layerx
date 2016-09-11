package framework_api_handlers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/layerx_rpi_client"
	"github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/mesos/mesos-go/mesosproto"
)

func HandleResourceOffers(lxRpi *layerx_rpi_client.LayerXRpi, offers []*mesosproto.Offer) error {
	for _, offer := range offers {
		resource := lxtypes.NewResourceFromMesos(offer)
		err := lxRpi.SubmitResource(resource)
		if err != nil {
			return errors.New("failed to submit resource "+resource.Id+" to layerx core", err)
		}
	}
	return nil
}
