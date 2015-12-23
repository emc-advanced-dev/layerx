package driver
import (
	"github.com/layer-x/layerx-commons/lxactionqueue"
)

type tpiDriver struct {
	actionQueue lxactionqueue.ActionQueue
}

func NewMesosTpiDriver(actionQueue lxactionqueue.ActionQueue) *tpiDriver {
	return &tpiDriver{
		actionQueue: actionQueue,
	}
}

//run as goroutine
func (d *tpiDriver) Run() {
	for {
		d.actionQueue.ExecuteNext()
	}
}