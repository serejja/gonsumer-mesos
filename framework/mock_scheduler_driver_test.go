package framework

import mesos "github.com/mesos/mesos-go/mesosproto"

type MockSchedulerDriver struct {
	StartStatus mesos.Status
	StartError  error

	StopStatus mesos.Status
	StopError  error

	AbortStatus mesos.Status
	AbortError  error

	JoinStatus mesos.Status
	JoinError  error

	RunStatus mesos.Status
	RunError  error

	RequestResourcesStatus mesos.Status
	RequestResourcesError  error

	LaunchTasksStatus mesos.Status
	LaunchTasksError  error

	KillTaskStatus mesos.Status
	KillTaskError  error
	KillTaskCount  int

	DeclineOfferStatus mesos.Status
	DeclineOfferError  error

	ReviveOffersStatus mesos.Status
	ReviveOffersError  error

	SendFrameworkMessageStatus mesos.Status
	SendFrameworkMessageError  error

	ReconcileTasksStatus mesos.Status
	ReconcileTasksError  error
	ReconcileTasksCount  int
}

func NewMockSchedulerDriver() *MockSchedulerDriver {
	return &MockSchedulerDriver{
		StartStatus:                mesos.Status_DRIVER_RUNNING,
		StopStatus:                 mesos.Status_DRIVER_RUNNING,
		AbortStatus:                mesos.Status_DRIVER_RUNNING,
		JoinStatus:                 mesos.Status_DRIVER_RUNNING,
		RunStatus:                  mesos.Status_DRIVER_RUNNING,
		RequestResourcesStatus:     mesos.Status_DRIVER_RUNNING,
		LaunchTasksStatus:          mesos.Status_DRIVER_RUNNING,
		KillTaskStatus:             mesos.Status_DRIVER_RUNNING,
		DeclineOfferStatus:         mesos.Status_DRIVER_RUNNING,
		ReviveOffersStatus:         mesos.Status_DRIVER_RUNNING,
		SendFrameworkMessageStatus: mesos.Status_DRIVER_RUNNING,
		ReconcileTasksStatus:       mesos.Status_DRIVER_RUNNING,
	}
}

func (s *MockSchedulerDriver) Start() (mesos.Status, error) {
	return s.StartStatus, s.StartError
}

func (s *MockSchedulerDriver) Stop(failover bool) (mesos.Status, error) {
	return s.StopStatus, s.StopError
}

func (s *MockSchedulerDriver) Abort() (mesos.Status, error) {
	return s.AbortStatus, s.AbortError
}

func (s *MockSchedulerDriver) Join() (mesos.Status, error) {
	return s.JoinStatus, s.JoinError
}

func (s *MockSchedulerDriver) Run() (mesos.Status, error) {
	return s.RunStatus, s.RunError
}

func (s *MockSchedulerDriver) RequestResources(requests []*mesos.Request) (mesos.Status, error) {
	return s.RequestResourcesStatus, s.RequestResourcesError
}

func (s *MockSchedulerDriver) LaunchTasks(offerIDs []*mesos.OfferID, tasks []*mesos.TaskInfo, filters *mesos.Filters) (mesos.Status, error) {
	return s.LaunchTasksStatus, s.LaunchTasksError
}

func (s *MockSchedulerDriver) KillTask(taskID *mesos.TaskID) (mesos.Status, error) {
	s.KillTaskCount++
	return s.KillTaskStatus, s.KillTaskError
}

func (s *MockSchedulerDriver) DeclineOffer(offerID *mesos.OfferID, filters *mesos.Filters) (mesos.Status, error) {
	return s.DeclineOfferStatus, s.DeclineOfferError
}

func (s *MockSchedulerDriver) ReviveOffers() (mesos.Status, error) {
	return s.ReviveOffersStatus, s.ReviveOffersError
}

func (s *MockSchedulerDriver) SendFrameworkMessage(executorID *mesos.ExecutorID, slaveID *mesos.SlaveID, data string) (mesos.Status, error) {
	return s.SendFrameworkMessageStatus, s.SendFrameworkMessageError
}

func (s *MockSchedulerDriver) ReconcileTasks(statuses []*mesos.TaskStatus) (mesos.Status, error) {
	s.ReconcileTasksCount++
	return s.ReconcileTasksStatus, s.ReconcileTasksError
}
