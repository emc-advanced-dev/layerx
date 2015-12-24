package driver
import (
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"sync"
)

type tpiDriver struct {
	lock sync.Mutex
	actionQueue lxactionqueue.ActionQueue
}

func NewMesosTpiDriver(actionQueue lxactionqueue.ActionQueue) *tpiDriver {
	return &tpiDriver{
		lock: sync.Mutex{},
		actionQueue: actionQueue,
	}
}

//run as goroutine
func (d *tpiDriver) Run() {
	for {
		d.lock.Lock()
		d.actionQueue.ExecuteNext()
		d.lock.Unlock()
	}
}