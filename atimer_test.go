// Written by Gon Yi
package atimer_test

import (
	"github.com/orangenumber/atimer"
	"testing"
	"time"
)

func TestNewJob(t *testing.T) {
	j1 := atimer.NewJob(func(j *atimer.Job) {
		println(time.Now().String(), "every second", j.GetName())
	}, time.Second)

	j2 := atimer.NewJob(func(j *atimer.Job) {
		println(time.Now().String(), "every 2 second", j.GetName())
	}, time.Second*2)

	for {
		j1.Do()
		j2.Do()
	}
}

func TestNew(t *testing.T) {
	timer := atimer.NewBucket()
	j1 := timer.NewJob(func(j *atimer.Job) {
		println(time.Now().String(), "every sec", j.GetName())
	}, time.Second)
	timer.Start(time.Second)

	time.Sleep(time.Second * 3)
	j2 := atimer.NewJob(func(j *atimer.Job) {
		println(time.Now().String(), "every 2 second", j.GetName())
	}, time.Second*2)
	timer.Link(j2)

	time.Sleep(time.Second * 4)
	println("name", j2.GetName())
	j2.Rename("testCounter")
	println("name", j2.GetName())

	j2.SetInterval(time.Second)
	j2.SetDo(func(j *atimer.Job) {
		println("now every sec!")
	})
	time.Sleep(2 * time.Second)
	println("Remove original sec counter")
	j1.Unregister()

	time.Sleep(10 * time.Second)
	timer.Stop()
	println("stops")
	time.Sleep(3 * time.Second)
}

func TestNewJobRun(t *testing.T) {
	j := atimer.NewJob(func(j *atimer.Job) {
		println(time.Now().Format("03:04:05"), "every half sec")
	}, time.Second)
	j.Start(0)
	time.Sleep(time.Second * 5)
	j.Stop()
	time.Sleep(time.Second * 5)
	j.Start(time.Second * 2)
	time.Sleep(time.Second * 5)
}

func TestNewJobRun2(t *testing.T) {
	// Run a single and then include in the bucket
	j := atimer.NewJob(func(j *atimer.Job) {
		println(time.Now().Format("03:04:05"), "every half sec")
	}, time.Second)
	j.Start(0)
	time.Sleep(time.Second * 3)

	b := atimer.NewBucket()
	b.Link(j)

	time.Sleep(time.Second * 3)

	b.Start(time.Second)
	time.Sleep(time.Second * 3)
}
