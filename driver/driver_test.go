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
		It("asynchronously executes functions in the action queue", func(){
			resultchan := make(chan int)
			fun1 := func(){
				resultchan <- 1
			}

			fun2 := func(){
				resultchan <- 2
			}

			fun3 := func(){
				resultchan <- 3
			}

			actionQueue.Push(fun1)
			actionQueue.Push(fun2)
			actionQueue.Push(fun3)

			result1 := <-resultchan
			result2 := <- resultchan
			result3 := <- resultchan
			results := []int{result1, result2, result3}
			Expect(results).To(ContainElement(1))
			Expect(results).To(ContainElement(2))
			Expect(results).To(ContainElement(3))
		})
	})
})
