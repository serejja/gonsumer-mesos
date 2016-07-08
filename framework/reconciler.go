package framework

import (
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	"github.com/mesos/mesos-go/scheduler"
	"github.com/yanzay/log"
	"sync"
	"time"
)

type Reconciler struct {
	ReconcileDelay    time.Duration
	ReconcileMaxTries int

	tasks         map[string]struct{}
	taskLock      sync.Mutex
	reconcileTime time.Time
	reconciles    int
}

func NewReconciler() *Reconciler {
	return &Reconciler{
		ReconcileDelay:    10 * time.Second,
		ReconcileMaxTries: 3,
		tasks:             make(map[string]struct{}),
		reconcileTime:     time.Unix(0, 0),
	}
}

func (r *Reconciler) ImplicitReconcile(driver scheduler.SchedulerDriver) error {
	return r.reconcile(driver, true)
}

func (r *Reconciler) ExplicitReconcile(taskIDs []string, driver scheduler.SchedulerDriver) error {
	r.taskLock.Lock()
	for _, taskID := range taskIDs {
		r.tasks[taskID] = struct{}{}
	}
	r.taskLock.Unlock()

	return r.reconcile(driver, false)
}

func (r *Reconciler) Update(status *mesos.TaskStatus) {
	r.taskLock.Lock()
	defer r.taskLock.Unlock()

	delete(r.tasks, status.GetTaskId().GetValue())

	if len(r.tasks) == 0 {
		r.reconciles = 0
	}
}

func (r *Reconciler) reconcile(driver scheduler.SchedulerDriver, implicit bool) error {
	if time.Now().Sub(r.reconcileTime) >= r.ReconcileDelay {
		r.taskLock.Lock()
		defer r.taskLock.Unlock()

		r.reconciles++
		r.reconcileTime = time.Now()

		if r.reconciles > r.ReconcileMaxTries {
			for task := range r.tasks {
				log.Infof("Reconciling exceeded %d tries, sending killTask for task %s", r.ReconcileMaxTries, task)
				_, err := driver.KillTask(util.NewTaskID(task))
				if err != nil {
					return err
				}
			}
			r.reconciles = 0
		} else {
			if implicit {
				_, err := driver.ReconcileTasks(nil)
				if err != nil {
					return err
				}
			} else {
				statuses := make([]*mesos.TaskStatus, 0)
				for task := range r.tasks {
					log.Debugf("Reconciling %d/%d task state for task id %s", r.reconciles, r.ReconcileMaxTries, task)
					statuses = append(statuses, util.NewTaskStatus(util.NewTaskID(task), mesos.TaskState_TASK_STAGING))
				}
				_, err := driver.ReconcileTasks(statuses)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
