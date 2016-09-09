package fakes

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/mesosproto/scheduler"
)

func FakeSubscribeCall() *scheduler.Call {
	callType := scheduler.Call_SUBSCRIBE
	return &scheduler.Call{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String("fake_framework_id"),
		},
		Type: &callType,
		Subscribe: &scheduler.Call_Subscribe{
			FrameworkInfo: FakeFramework(),
		},
	}
}

func FakeDeclineOffersCall(frameworkId string, offerIds ...string) *scheduler.Call {
	callType := scheduler.Call_DECLINE
	mesosOfferIds := []*mesosproto.OfferID{}
	for _, offerId := range offerIds {
		mesosOfferIds = append(mesosOfferIds, &mesosproto.OfferID{
			Value: proto.String(offerId),
		})
	}
	return &scheduler.Call{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String(frameworkId),
		},
		Type: &callType,
		Decline: &scheduler.Call_Decline{
			OfferIds: mesosOfferIds,
		},
	}
}

func FakeReconcileTasksCall(frameworkId string, taskIds ...string) *scheduler.Call {
	callType := scheduler.Call_RECONCILE
	reconcileTasks := []*scheduler.Call_Reconcile_Task{}
	for _, taskId := range taskIds {
		reconcileTasks = append(reconcileTasks, &scheduler.Call_Reconcile_Task{
			TaskId: &mesosproto.TaskID{
				Value: proto.String(taskId),
			},
		})
	}
	return &scheduler.Call{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String(frameworkId),
		},
		Type: &callType,
		Reconcile: &scheduler.Call_Reconcile{
			Tasks: reconcileTasks,
		},
	}
}

func FakeReviveOffersCall(frameworkId string) *scheduler.Call {
	callType := scheduler.Call_REVIVE
	return &scheduler.Call{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String(frameworkId),
		},
		Type: &callType,
	}
}

func FakeLaunchTasksCall(frameworkId string, offerIds []string, taskInfos ...*mesosproto.TaskInfo) *scheduler.Call {
	callType := scheduler.Call_ACCEPT
	mesosOfferIds := []*mesosproto.OfferID{}
	for _, offerId := range offerIds {
		mesosOfferIds = append(mesosOfferIds, &mesosproto.OfferID{
			Value: proto.String(offerId),
		})
	}
	operationType := mesosproto.Offer_Operation_LAUNCH
	launchOperations := []*mesosproto.Offer_Operation{}
	for _, taskInfo := range taskInfos {
		launchOperations = append(launchOperations, &mesosproto.Offer_Operation{
			Type: &operationType,
			Launch: &mesosproto.Offer_Operation_Launch{
				TaskInfos: []*mesosproto.TaskInfo{taskInfo},
			},
		})
	}
	return &scheduler.Call{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String(frameworkId),
		},
		Type: &callType,
		Accept: &scheduler.Call_Accept{
			OfferIds:   mesosOfferIds,
			Operations: launchOperations,
		},
	}
}

func FakeKillTaskCall(frameworkId, taskId string) *scheduler.Call {
	callType := scheduler.Call_KILL
	return &scheduler.Call{
		FrameworkId: &mesosproto.FrameworkID{
			Value: proto.String(frameworkId),
		},
		Type: &callType,
		Kill: &scheduler.Call_Kill{
			TaskId: &mesosproto.TaskID{
				Value: proto.String(taskId),
			},
		},
	}
}
