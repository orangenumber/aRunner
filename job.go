// Written by Gon Yi
package atimer

import (
	"errors"
	"time"
)

type Job struct {
	name                 string
	parent               *Bucket
	do                   func(job *Job)
	lastRun              time.Time
	interval             time.Duration
	nextRun              time.Time
	enable               bool
	runningWithoutBucket bool
}

func newJob(parent *Bucket, name string, do func(job *Job), interval time.Duration) *Job {
	j := Job{
		name:     name,
		parent:   parent,
		do:       do,
		interval: interval,
		enable:   true,
	}

	// if parent is not nil, register
	if parent != nil {
		j.Register(parent)
	}

	return &j
}

func NewJob(do func(job *Job), interval time.Duration) *Job {
	return newJob(nil, "", do, interval)
}

//
// SetDo
//
func (j *Job) SetDo(do func(job *Job)) {
	j.do = do
}

//
// SetInterval will change the interval; and set nextRun time.
//
func (j *Job) SetInterval(t time.Duration) {
	j.interval = t
	j.nextRun = j.lastRun.Add(j.interval)
}

//
// SetNextRun will reset interval to -1, therefore it won't run over and over.
//
func (j *Job) SetNextRun(t time.Time) {
	j.interval = -1
	j.nextRun = t
}

//
// GetNextRun will return nextRun time.
//
func (j *Job) GetNextRun() time.Time {
	return j.nextRun
}

//
// Enable the Job
//
func (j *Job) Enable() {
	j.enable = true
}

//
// Disable the Job -- not removing completely
//
func (j *Job) Disable() {
	j.enable = false
}

//
// Register to Bucket
//
func (j *Job) Register(a *Bucket) error {
	newName := j.name
	if j.name == "" { // otherwise, random string to be used for its name
		foundEmpty := false
		for foundEmpty == false {
			tmp := "atimer." + randChars(20)
			if _, ok := a.jobs[tmp]; !ok { // if not exist, then use it.
				newName = tmp
				foundEmpty = true
			}
		}
	}
	if a == nil {
		return errors.New("param *Bucket cannot be nil")
	}
	j.parent = a
	j.name = newName
	j.runningWithoutBucket = false // if a job was running itself without a bucket, disable it.
	a.jobs[j.name] = j

	return nil
}

//
// Unregister from Bucket
//
func (j *Job) Unregister() {
	// If parent exist, remove itself
	if j.parent != nil && j.name != "" {
		delete(j.parent.jobs, j.name)
	}
}

//
// Rename the job
//
func (j *Job) Rename(name string) {
	j.Unregister()
	j.name = name
	j.Register(j.parent)
}

//
// GetName returns the name of the job
//
func (j *Job) GetName() string {
	return j.name
}

//
// ForceDo -- run the function and update lastRun time WITHOUT time or enabled.
//
func (j *Job) ForceDo() {
	j.do(j)
	j.lastRun = time.Now()
	if j.interval > 0 { // only schedule nextRun if interval is set.
		j.nextRun = j.lastRun.Add(j.interval)
	}
}

//
// Do -- run the function if ENABLED and DUE
// 	  returns true if the func ran, otherwise false.
//
func (j *Job) Do() bool {
	// if enabled, AND expected next run is due, run it.
	if j.enable == true && time.Now().After(j.nextRun) {
		j.ForceDo()
		return true
	}
	return false
}

//
// Start -- running itself without a bucket
//
func (j *Job) Start(refreshInterval time.Duration) {
	j.runningWithoutBucket = true
	go func() {
		for j.runningWithoutBucket {
			if j.runningWithoutBucket == true {
				j.Do()
				time.Sleep(refreshInterval)
			} else {
				break
			}
		}
	}()
}

//
// Stop -- stop running
//
func (j *Job) Stop() {
	j.runningWithoutBucket = false
}
