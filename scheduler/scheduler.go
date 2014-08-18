// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package scheduler is responsable for executing jobs periodically
package scheduler

import (
	"errors"
	"sync"
	"time"
)

var (
	// Number of seconds that the scheduler will before it check again all jobs
	SchedulerExecutionInterval time.Duration = 60 * time.Second
)

var (
	ErrJobTypeNotFound = errors.New("Job type not found in scheduler")
)

var (
	jobsMutex sync.Mutex // Lock to make it possible to add jobs after the scheduler already started
	jobs      []Job      // List of jobs that are going to be executed
)

// List of possible job types on the scheduler
const (
	JobTypeUnknown      JobType = 0 // When the job is not going to be verified later
	JobTypeScan         JobType = 1 // We usually wants to known when the next scan will start
	JobTypeNotification JobType = 2 // Identify notification jobs
)

// Identify the jobs with types to later retrieve information from scheduler to known when
// a specific job will run
type JobType int

// Job struct store all necessary information to execute a task periodically in the
// system. You can define a specific execution time and the interval that it will be
// executed
type Job struct {
	Type          JobType       // Type of the object to be identified later
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
		job.NextExecution = time.Now().UTC().Add(job.Interval)
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
					if time.Now().UTC().After(job.NextExecution) {
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

// Retrieve the next execution of the first occurance of a specific job type. If not found
// an error is returned
func NextExecutionByType(jobType JobType) (time.Time, error) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	for _, job := range jobs {
		if job.Type == jobType {
			return job.NextExecution, nil
		}
	}

	return time.Time{}, ErrJobTypeNotFound
}
