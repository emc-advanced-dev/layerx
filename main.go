package main
import (
	"github.com/layer-x/layerx-commons/lxlog"
	"github.com/Sirupsen/logrus"
	"github.com/layer-x/layerx-core_v2/fakes"
	"github.com/mesos/mesos-go/mesosproto"
)

func main(){
	fakeStatus1 := fakes.FakeTaskStatus("fake_task_id_1", mesosproto.TaskState_TASK_RUNNING)
	fakeStatus2 := fakes.FakeTaskStatus("fake_task_id_2", mesosproto.TaskState_TASK_KILLED)
	fakeStatus3 := fakes.FakeTaskStatus("fake_task_id_3", mesosproto.TaskState_TASK_FINISHED)

	fakeStatuses := []*mesosproto.TaskStatus{fakeStatus1, fakeStatus2, fakeStatus3}

	lxlog.Infof(logrus.Fields{"port": 6666}, "running FAKE layerx!")
	fakes.RunFakeLayerXServer(fakeStatuses, 6666)
}