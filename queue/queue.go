package queue

import (
	"runtime"
	"time"
)

//Clock ...
type Clock struct {
	timer time.Duration

	queue     chan bool
	queueDone chan bool
}

//NewClock ...
func NewClock(timer time.Duration, lenQueue uint) (c *Clock) {
	c = new(Clock)

	c.timer = timer
	c.queue = make(chan bool, lenQueue)
	c.queueDone = make(chan bool, lenQueue)

	go func() {
		l := time.Tick(c.timer)
		for range l {
			lenQ := len(c.queueDone)

			for i := 0; i < lenQ; i++ {
				<-c.queue
				<-c.queueDone
			}
		}
	}()

	return
}

//Add ...
func (c *Clock) Add(count uint) {
	var i uint
	for i = 0; i < count; i++ {
		c.queue <- true
	}
}

//Done ...
func (c *Clock) Done(count uint) {
	var i uint
	for i = 0; i < count; i++ {
		c.queueDone <- true
	}
}

//Wait ...
func (c *Clock) Wait() {
	runtime.Gosched()

	c.queue <- true
	<-c.queue
}

//WaitBool ...
func (c Clock) WaitBool() (wait bool) {
	runtime.Gosched()

	if len(c.queue) == cap(c.queue) {
		wait = true
	}

	return
}
