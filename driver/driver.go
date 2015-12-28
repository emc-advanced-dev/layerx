package driver
import (
	"github.com/layer-x/layerx-commons/lxactionqueue"
	"sync"
)

type rpiDriver struct {
	lock sync.Mutex
	actionQueue lxactionqueue.ActionQueue
}

func NewMesosRpiDriver(actionQueue lxactionqueue.ActionQueue) *rpiDriver {
	return &rpiDriver{
		lock: sync.Mutex{},
		actionQueue: actionQueue,
	}
}

//run as goroutine
func (d *rpiDriver) Run() {
	for {
		d.lock.Lock()
		go d.actionQueue.ExecuteNext()
		d.lock.Unlock()
	}
}