package driver
import (
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"sync"
)

type coreDriver struct {
	lock sync.Mutex
	actionQueue lxactionqueue.ActionQueue
}

func NewLayerXDriver(actionQueue lxactionqueue.ActionQueue) *coreDriver {
	return &coreDriver{
		lock: sync.Mutex{},
		actionQueue: actionQueue,
	}
}

//run as goroutine
func (d *coreDriver) Run() {
	for {
		func() {
			d.lock.Lock()
			defer d.lock.Unlock()
			d.actionQueue.ExecuteNext()
		}()
	}
}