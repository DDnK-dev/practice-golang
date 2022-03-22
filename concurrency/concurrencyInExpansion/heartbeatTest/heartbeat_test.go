package main

import (
	"testing"
	"time"
)

func TestProducerFixedInterval(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 2 * time.Second
	heartbeat, results := ProducerFixedInterval(done, timeout/2, intSlice...)

	<-heartbeat

	i := 0
	for {
		select {
		case r, ok := <-results:
			if ok == false {
				return
			} else if expected := intSlice[i]; r != expected {
				t.Errorf(
					"index %v: expected %v, but recieved %v",
					i,
					expected,
					r,
				)
			}
			i++
		case <-heartbeat: // 시간초과가 발생하지 않도록 여기서 하트비트를 select
		case <-time.After(timeout):
			t.Fatalf("test timed out")
		}
	}
}
