package framework

import (
	"errors"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestImplicitReconcile(t *testing.T) {
	reconciler := NewReconciler()
	driver := NewMockSchedulerDriver()

	timestamp := time.Now()

	err := reconciler.ImplicitReconcile(driver)
	assert.Nil(t, err)
	assert.Equal(t, 1, reconciler.reconciles)
	assert.True(t, reconciler.reconcileTime.After(timestamp))
	assert.Equal(t, 1, driver.ReconcileTasksCount)

	// second reconciliation call should be ignored due to ReconcileDelay
	err = reconciler.ImplicitReconcile(driver)
	assert.Nil(t, err)
	assert.Equal(t, 1, reconciler.reconciles)
	assert.Equal(t, 1, driver.ReconcileTasksCount)

	// reconciler should propagate driver errors
	reconciler = NewReconciler()
	driver.ReconcileTasksError = errors.New("boom!")
	err = reconciler.ImplicitReconcile(driver)
	assert.EqualError(t, err, "boom!")
}

func TestExplicitReconcile(t *testing.T) {
	reconciler := NewReconciler()
	driver := NewMockSchedulerDriver()

	timestamp := time.Now()

	err := reconciler.ExplicitReconcile([]string{"foo", "bar"}, driver)
	assert.Nil(t, err)
	assert.Equal(t, 1, reconciler.reconciles)
	assert.True(t, reconciler.reconcileTime.After(timestamp))
	assert.Equal(t, 1, driver.ReconcileTasksCount)
	assert.Len(t, reconciler.tasks, 2)

	// reconciler should propagate driver errors
	reconciler = NewReconciler()
	driver.ReconcileTasksError = errors.New("boom!")
	err = reconciler.ExplicitReconcile([]string{"foo", "bar"}, driver)
	assert.EqualError(t, err, "boom!")

	// reconciler should kill tasks if MaxTries exceeded
	driver.ReconcileTasksError = nil
	reconciler.ReconcileDelay = time.Duration(0)
	reconciler.ReconcileMaxTries = 0

	err = reconciler.ExplicitReconcile([]string{"foo", "bar"}, driver)
	assert.Nil(t, err)
	assert.Equal(t, 0, reconciler.reconciles) // should reset
	assert.Equal(t, 2, driver.KillTaskCount)

	// reconciler should propagate driver errors if fails to kill a task
	driver.KillTaskError = errors.New("boom!")
	err = reconciler.ExplicitReconcile([]string{"foo", "bar"}, driver)
	assert.EqualError(t, err, "boom!")
}

func TestReconcilerUpdate(t *testing.T) {
	reconciler := NewReconciler()
	reconciler.tasks["foo"] = struct{}{}
	reconciler.tasks["bar"] = struct{}{}
	reconciler.reconciles = 2

	reconciler.Update(&mesos.TaskStatus{
		TaskId: util.NewTaskID("foo"),
	})

	assert.Len(t, reconciler.tasks, 1)
	assert.Equal(t, 2, reconciler.reconciles)

	reconciler.Update(&mesos.TaskStatus{
		TaskId: util.NewTaskID("bar"),
	})

	assert.Len(t, reconciler.tasks, 0)
	assert.Equal(t, 0, reconciler.reconciles) // reset
}
