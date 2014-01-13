package scheduler

import (
	"sync"
	"time"
)

const (
	// Number of minutes that the scheduler will before it check again all jobs
	SchedulerExecutionInterval = 1
)

var (
	jobsMutex sync.Mutex // Lock to make it possible to add jobs after the scheduler already started
	jobs      []Job      // List of jobs that are going to be executed
)

func init() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				jobsMutex.Lock()
				for index, job := range jobs {
					if time.Now().After(job.NextExecution) {
						jobs[index].NextExecution = time.Now().Add(job.Interval)
						go job.Task()
					}
				}
				jobsMutex.Unlock()
			}
		}
	}()
}

type Job struct {
	NextExecution time.Time
	Interval      time.Duration
	Task          func()
}

func Register(job Job) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()
	jobs = append(jobs, job)
}
