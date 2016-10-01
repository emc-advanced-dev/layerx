package lx_core_helpers

import (
	"github.com/emc-advanced-dev/layerx/layerx-core/lxstate"
	"github.com/emc-advanced-dev/pkg/errors"
)

func RescindResource(state *lxstate.State, resourceID string) error {
	resourcePool, err := state.GetResourcePoolContainingResource(resourceID)
	if err != nil {
		return errors.New("failed to get resource pool for resource", err)
	}
	if err := resourcePool.DeleteResource(resourceID); err != nil {
		return errors.New("removing resource from pool", err)
	}
	return nil
}
