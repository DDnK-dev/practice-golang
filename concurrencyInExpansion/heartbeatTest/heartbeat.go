package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	//consumer()// consume code
	//consumer2()
	consumer3()
}

// producerForTest 작업 단위의 시작부분에서 발생하는 하트비트를 살펴보는 예제이다.
// 테스트에 유용하다...
func consumer3() {
	done := make(chan interface{})
	defer close(done)

	heartbeat, results := producerForTest(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v \n", r)
			} else {
				return
			}
		}
	}
}

func producerForTest(done <-chan interface{}) (<-chan interface{}, <-chan int) {
	// 크기가 1인 버퍼로 생성. 송신 대기 시간 내에 아무도 채널을 듣고 있지 않아도 적어도 하나의 펄스가 송출
	heartbeatStream := make(chan interface{}, 1)
	workStream := make(chan int)
	go func() {
		defer close(heartbeatStream)
		defer close(workStream)

		for i := 0; i < 10; i++ {
			// heartbeat를 위한 별도의 select 블록 설정
			// 수신자가 결과를 바등ㄹ 준비가 되지 않았다면 결과 대신 펄스를 받게 되고, 현재의 결과 값을 잃어버리게
			// 되므로 results 채널에 대한 전송과 동일한 블록에 포함 X.
			select {
			case heartbeatStream <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case workStream <- rand.Intn(10):
			}
		}
	}()
	return heartbeatStream, workStream
}

// 여기까지가 timeout을 확인하는 2번 예제
func consumer2() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := producerMakingPanic(done, timeout/2)

	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r)
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healty!")
			return
		}
	}
}

func producerMakingPanic(
	done <-chan interface{},
	pulseInterval time.Duration,
) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)
	go func() {
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2 * pulseInterval)
		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}
		sendResult := func(r time.Time) {
			for {
				select {
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		// 2회 실행 후에 열려있는 채널을 그대로 놔둔채로 끝내버린다..!!
		for i := 0; i < 2; i++ {
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

// 이 아래는 1번 예제
func consumer() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := producer(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
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
