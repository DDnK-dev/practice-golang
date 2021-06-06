package main

import "time"

func main() {
	done := make(chan interface{})
	producer(done, time.Second*1)
}

// 하트비트 채널을 만들고 이벤트를 발생시키는 함수
func producer(
	done <-chan interface{},
	pulseInterval time.Duration,
) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{}) // 하트비트를 보내기 위한 채널
	results := make(chan time.Time)
	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.Tick(pulseInterval)       // 주어진 pulseInterval에서 하트비트가 뛰도록 설정. 이때는 이 채널에서 읽을 내용이 있음
		workGen := time.Tick(2 * pulseInterval) // 들어오는 작업을 시뮬하는데 사용되는 또 다른 티커... 테스트용임

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default: // 언제나 아무도 하트비트를 듣지 않을 수 있다는 사실에 대비. 그것 자체는 별로 중요하지 않으니..
			}
		}
		sendResult := func(r time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse: // 송수신 할 때도 포함
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		// 실제 실행부
		for {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()

	return heartbeat, results
}

// 이벤트를 소비하는 함수
