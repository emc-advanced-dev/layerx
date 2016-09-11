package fakes

import (
	"fmt"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/go-martini/martini"
	"github.com/gogo/protobuf/proto"

	"github.com/emc-advanced-dev/layerx/layerx-mesos-tpi/mesos_master_api/mesos_data"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/mesos/mesos-go/mesosproto"
	"io/ioutil"
	"net/http"
)

func RunFakeFrameworkServer(frameworkid string, port int) {

	var offersRecieved = 0

	m := martini.Classic()
	m.Post("/"+frameworkid+"/mesos.internal.FrameworkRegisteredMessage", func(req *http.Request, res http.ResponseWriter) {
		body, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			res.Header().Add("error", err.Error())
			res.WriteHeader(500)
			return
		}
		var frameworkRegistered mesosproto.FrameworkRegisteredMessage
		err = proto.Unmarshal(body, &frameworkRegistered)
		if err != nil {
			fmt.Printf("\nerr: %v\n", err)
			res.Header().Add("error", errors.New("received data("+string(body)+")", err).Error())
			res.WriteHeader(500)
			return
		}
		fmt.Printf("finished")
		res.WriteHeader(202)
	})
	m.Post("/"+frameworkid+"/mesos.internal.ResourceOffersMessage", func(req *http.Request, res http.ResponseWriter) {
		body, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			res.Header().Add("error", err.Error())
			res.WriteHeader(500)
			return
		}
		var resourceOffers mesosproto.ResourceOffersMessage
		err = proto.Unmarshal(body, &resourceOffers)
		if err != nil {
			fmt.Printf("\nerr: %v\n", err)
			res.Header().Add("error", errors.New("received data("+string(body)+")", err).Error())
			res.WriteHeader(500)
			return
		}
		if len(resourceOffers.Offers) == 0 {
			fmt.Printf("\noffers recieved: %v\n", len(resourceOffers.Offers))
			res.Header().Add("error", errors.New("received only 0 offers", nil).Error())
			res.WriteHeader(500)
			return
		}
		masterPidString := req.Header.Get("Libprocess-From")
		if masterPidString == "" {
			fmt.Printf("missing required header: %s", "Libprocess-From")
			res.WriteHeader(400)
			return
		}
		upid, err := mesos_data.UPIDFromString(masterPidString)
		if err != nil {
			fmt.Printf("could not parse upid of master")
			res.WriteHeader(400)
			return
		}
		for _, offer := range resourceOffers.Offers {
			offersRecieved++
			sendFakeTaskOnOffer(offersRecieved, offer, upid, port)
		}

		fmt.Printf("finished")
		res.WriteHeader(202)
	})
	m.Post("/"+frameworkid+"/mesos.internal.StatusUpdateMessage", func(req *http.Request, res http.ResponseWriter) {
		body, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			res.Header().Add("error", err.Error())
			res.WriteHeader(500)
			return
		}
		var statusUpdate mesosproto.StatusUpdateMessage
		err = proto.Unmarshal(body, &statusUpdate)
		if err != nil {
			fmt.Printf("\nerr: %v\n", err)
			res.Header().Add("error", errors.New("received data("+string(body)+")", err).Error())
			res.WriteHeader(500)
			return
		}
		fmt.Printf("finished")
		res.WriteHeader(202)
	})

	m.RunOnAddr(fmt.Sprintf(":%v", port))
}

func sendFakeTaskOnOffer(taskNo int, offer *mesosproto.Offer, masterUpid *mesos_data.UPID, port int) {
	taskInfo := &mesosproto.TaskInfo{
		Name: proto.String(fmt.Sprintf("fake-task-%v", taskNo)),
		TaskId: &mesosproto.TaskID{
			Value: proto.String(fmt.Sprintf("fake-taskId-%v", taskNo)),
		},
		SlaveId: &mesosproto.SlaveID{
			Value: proto.String(offer.GetSlaveId().GetValue()),
		},
		Resources: offer.Resources,
	}
	fmt.Printf("\n\ntask i will send: %v\n", taskInfo.GetTaskId().GetValue())

	fakeTasks := []*mesosproto.TaskInfo{
		taskInfo,
	}
	fakeOfferIds := []*mesosproto.OfferID{
		&mesosproto.OfferID{
			Value: proto.String("fake_offer_id"),
		},
	}
	launch := &mesosproto.LaunchTasksMessage{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String("fake_framework_id"),
		},
		Tasks:    fakeTasks,
		Filters:  &mesosproto.Filters{},
		OfferIds: fakeOfferIds,
	}

	headers := map[string]string{
		"Libprocess-From": fmt.Sprintf("fakeframework@127.0.0.1:%v", port),
	}
	fmt.Printf("\nPOSTING TO LAUNCH TASKS on addr %s\n", masterUpid.Host+":"+masterUpid.Port)
	_, _, err := lxhttpclient.Post(masterUpid.Host+":"+masterUpid.Port,
		"/master/mesos.internal.LaunchTasksMessage",
		headers,
		launch)
	if err != nil {
		fmt.Printf("\nBIG ERROR POSTING LAUNCH TASKS!! WHAT'S THE DEAL!?\n%v\n", err)
	}
}
