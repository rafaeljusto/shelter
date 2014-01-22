package scheduler

import (
	"fmt"
	"testing"
	"time"
)

var (
	ValueToChange int
)

func TestFlexibleTimeJobExecution(t *testing.T) {
	SchedulerExecutionInterval = 50 * time.Millisecond

	Clear()

	ValueToChange = 0
	Register(Job{
		Interval: SchedulerExecutionInterval / 2,
		Task: func() {
			ValueToChange += 1
		},
	})

	Start()

	// Job will not execute exactly in time, but after this amount of time we expect two
	// executions of our job. We assume a duration of 1 milisecond to execute the job two
	// times
	time.Sleep((SchedulerExecutionInterval * 2) + 10*time.Millisecond)

	if ValueToChange != 2 {
		t.Error(fmt.Sprintf("Not executing a flexible time job properly. "+
			"Expected %d and got %d", 2, ValueToChange))
	}
}

func TestSpecificTimeJobExecution(t *testing.T) {
	SchedulerExecutionInterval = 50 * time.Millisecond

	Clear()

	ValueToChange = 0
	Register(Job{
		NextExecution: time.Now().Add(100 * time.Millisecond),
		Interval:      SchedulerExecutionInterval / 2,
		Task: func() {
			ValueToChange += 1
		},
	})

	Start()

	// Job will not execute exactly in time, but after this amount of time we expect one
	// execution of our job, because it was scheduled for later. We assume a duration of 1
	// milisecond to execute the job one time
	time.Sleep((SchedulerExecutionInterval * 2) + 1*time.Millisecond)

	if ValueToChange != 1 {
		t.Error(fmt.Sprintf("Not executing a time specific job properly. "+
			"Expected %d and got %d", 1, ValueToChange))
	}
}
