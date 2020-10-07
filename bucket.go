// Written by Gon Yi
package atimer

import (
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

type Bucket struct {
	jobs            map[string]*Job
	mu              sync.Mutex
	enable          bool
	refreshInterval time.Duration
	log             io.Writer
}

func (r *Bucket) Log(category, format string, a ...interface{}) {
	r.log.Write([]byte(fmt.Sprintf("["+category+"] "+format+"\n", a...)))
}

func NewBucket() *Bucket {
	return &Bucket{
		jobs: make(map[string]*Job),
		log:  ioutil.Discard,
	}
}

func (r *Bucket) NewJob(do func(job *Job), interval time.Duration) *Job {
	return newJob(r, "", do, interval)
}

func (r *Bucket) Link(job *Job) {
	job.Register(r)
}

func (r *Bucket) Unlink(job *Job) {
	job.Unregister()
}

func (r *Bucket) Start(refreshInterval time.Duration) {
	r.refreshInterval = refreshInterval
	r.Log("OK", "starting the runner (refresh interval: %s)", refreshInterval.String())
	r.enable = true

	go func() { // run in the background
		for r.enable == true {
			r.Log("RUN", "refresh started")
			for k, v := range r.jobs {
				if v.Do() == true { // only run if enabled
					r.Log("Run", "   %s", k)
				}
			}
			r.Log("RUN", "refresh ended")
			time.Sleep(r.refreshInterval)
		}
	}()
}

func (r *Bucket) Stop() {
	r.Log("STOP", "stop the runner")
	r.enable = false
}
