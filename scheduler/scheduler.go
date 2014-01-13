package scheduler

import (
	"sync"
	"time"
)

const (
	SchedulerExecutionInterval = 1 // Number that defines the inetrval that the scheduler will be executed
)

var (
	jobsMutex sync.Mutex
	jobs      []Job
)

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

func Start() {
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
