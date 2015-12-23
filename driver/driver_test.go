package driver_test

import (
	. "github.com/layer-x/layerx-mesos-tpi_v2/driver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-commons/lxactionqueue"
)

var _ = Describe("Driver", func() {

	actionQueue := lxactionqueue.NewActionQueue()
	driver := NewMesosTpiDriver(actionQueue)

	go driver.Run()

	Describe("Driver main thread", func(){
		It("synchronously executes functions in the action queue", func(){
			resultchan := make(chan int)
			fun1 := func(){
				resultchan <- 1
			}
			actionQueue.Push(fun1)

			fun2 := func(){
				resultchan <- 2
			}
			actionQueue.Push(fun2)

			fun3 := func(){
				resultchan <- 3
			}
			actionQueue.Push(fun3)

			result1 := <-resultchan
			Expect(result1).To(Equal(1))
			result2 := <- resultchan
			Expect(result2).To(Equal(2))
			result3 := <- resultchan
			Expect(result3).To(Equal(3))
		})
	})
})
