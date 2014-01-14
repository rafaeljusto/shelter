package scheduler

import (
	"sync"
	"time"
)

var (
	// Number of seconds that the scheduler will before it check again all jobs
	SchedulerExecutionInterval time.Duration = 60 * time.Second
)

var (
	jobsMutex sync.Mutex // Lock to make it possible to add jobs after the scheduler already started
	jobs      []Job      // List of jobs that are going to be executed
)

// Job struct store all necessary information to execute a task periodically in the
// system. You can define a specific execution time and the interval that it will be
// executed
type Job struct {
	NextExecution time.Time     // Schedule the next execution
	Interval      time.Duration // Interval of executions of this job
	Task          func()        // Function that will be executed
}

// Function to register a new job, we use this instead of a global variable because we
// want to set a lock before changing the job list
func Register(job Job) {
	if job.NextExecution.IsZero() {
		// If the job does not have a next execution defined we assume that it does
		// not have an exactly time to run, so we just assume now as a reference
		job.NextExecution = time.Now().Add(job.Interval)
	}

	jobsMutex.Lock()
	defer jobsMutex.Unlock()
	jobs = append(jobs, job)
}

// Clear function was created for tests, so we can work on many scenarios without
// initializing the scheduler again and again
func Clear() {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()
	jobs = []Job{}
}

// We don't start scheduler in init function anymore because for tests we want to change
// parameters before the scheduler starts
func Start() {
	ticker := time.NewTicker(SchedulerExecutionInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				jobsMutex.Lock()
				for index, job := range jobs {
					// The execution time is not going to be so exactly
					if time.Now().After(job.NextExecution) {
						// Next execution time is defined, let's use it as reference so that the job
						// is always executed near the desired time
						jobs[index].NextExecution = job.NextExecution.Add(job.Interval)
						go job.Task()
					}
				}
				jobsMutex.Unlock()
			}
		}
	}()
}
