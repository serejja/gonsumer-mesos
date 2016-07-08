package framework

import (
	"encoding/json"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
	"github.com/serejja/gonsumer-mesos/mesosfmt"
	"github.com/yanzay/log"
	"time"
)

type Scheduler interface {
	Cluster() Cluster
}

type GonsumerScheduler struct {
	driver     scheduler.SchedulerDriver
	cluster    Cluster
	storage    Storage
	reconciler *Reconciler
}

func NewScheduler(storage Storage) (*GonsumerScheduler, error) {
	gonsumerScheduler := &GonsumerScheduler{
		storage:    storage,
		reconciler: NewReconciler(),
	}
	gonsumerScheduler.reconciler.ReconcileDelay = 30 * time.Second

	err := gonsumerScheduler.LoadClusterState()
	return gonsumerScheduler, err
}

func (s *GonsumerScheduler) Registered(driver scheduler.SchedulerDriver, id *mesos.FrameworkID, master *mesos.MasterInfo) {
	log.Infof("[Registered] framework: %s master: %s:%d", id.GetValue(), master.GetHostname(), master.GetPort())

	s.cluster.SetFrameworkID(id.GetValue())
	s.SaveClusterState()

	s.driver = driver
	s.reconciler.ImplicitReconcile(driver)
}

func (s *GonsumerScheduler) Reregistered(driver scheduler.SchedulerDriver, master *mesos.MasterInfo) {
	log.Infof("[Reregistered] master: %s:%d", master.GetHostname(), master.GetPort())

	s.driver = driver
	s.reconciler.ImplicitReconcile(driver)
}

func (s *GonsumerScheduler) Disconnected(scheduler.SchedulerDriver) {
	log.Info("[Disconnected]")
}

func (s *GonsumerScheduler) ResourceOffers(driver scheduler.SchedulerDriver, offers []*mesos.Offer) {
	log.Debugf("[ResourceOffers] %s", mesosfmt.Offers(offers))

	s.SaveClusterState() //TODO this should be only called when cluster state changed
}

func (s *GonsumerScheduler) StatusUpdate(driver scheduler.SchedulerDriver, status *mesos.TaskStatus) {
	log.Infof("[StatusUpdate] %s", mesosfmt.Status(status))

	s.SaveClusterState()
}

func (s *GonsumerScheduler) OfferRescinded(driver scheduler.SchedulerDriver, id *mesos.OfferID) {
	log.Infof("[OfferRescinded] %s", id.GetValue())
}

func (s *GonsumerScheduler) FrameworkMessage(driver scheduler.SchedulerDriver, executor *mesos.ExecutorID, slave *mesos.SlaveID, message string) {
	log.Infof("[FrameworkMessage] executor: %s slave: %s message: %s", executor.GetValue(), slave.GetValue(), message)
}

func (s *GonsumerScheduler) SlaveLost(driver scheduler.SchedulerDriver, slave *mesos.SlaveID) {
	log.Infof("[SlaveLost] %s", slave.GetValue())
}

func (s *GonsumerScheduler) ExecutorLost(driver scheduler.SchedulerDriver, executor *mesos.ExecutorID, slave *mesos.SlaveID, status int) {
	log.Infof("[ExecutorLost] executor: %s slave: %s status: %d", executor.GetValue(), slave.GetValue(), status)
}

func (s *GonsumerScheduler) Error(driver scheduler.SchedulerDriver, err string) {
	log.Errorf("[Error] %s", err)
}

func (s *GonsumerScheduler) Shutdown(driver scheduler.SchedulerDriver) {
	log.Info("Shutdown triggered, stopping driver")

	_, err := driver.Stop(false)
	if err != nil {
		panic(err)
	}
}

func (s *GonsumerScheduler) Cluster() Cluster {
	return s.cluster
}

func (s *GonsumerScheduler) LoadClusterState() error {
	rawCluster, err := s.storage.Load()
	if err == ErrStorageUninitialized {
		s.cluster = NewGonsumerCluster()
		return nil
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(rawCluster, &s.cluster)
}

func (s *GonsumerScheduler) SaveClusterState() error {
	clusterJSON, err := json.Marshal(s.cluster)
	if err != nil {
		return err
	}

	return s.storage.Save(clusterJSON)
}
