package main

//
//import (
//	"log"
//	"time"
//
//	"practice-golang/common"
//)
//
//// 스튜어드-와드 패턴을 만들고 비정상 고루틴을 감지하여 재시작(Healing) 하는 패턴을 보여준다
//// 하트비트 패턴을 이용한다.
//
//func main() {
//
//}
//
////startGoroutineFn 은 모니터링되고 재시작 될 고루틴의 시그니처를 정의한다. 하트비트가 사용됨을 볼 수 있다.
//type startGoroutineFn func(
//	done <-chan interface{},
//	pulseInterval time.Duration,
//) (heartbeat <-chan interface{})
//
//func newSteward(
//	timeout time.Duration,
//	startGoroutine startGoroutineFn,
//) startGoroutineFn {
//	// 모니터링 대상이 되는 고루틴에 대한 timeout, startGoroutine을 함수 인자로 받는 것을 볼 수 있다.
//	// startGoroutineFn을 리턴하는데, 이 자체도 모니터링이 가능하다는 것을 의미한다.
//	return func(
//		done <-chan interface{},
//		pulseInterval time.Duration,
//	) <-chan interface{} {
//		heartbeat := make(chan interface{})
//		go func() {
//			defer close(heartbeat)
//
//			var wardDone chan interface{}
//			var wardHeartbeat <-chan interface{}
//			startWard := func() { // 정형화된 방식으로 모니터링할 고루틴을 시작시킬 클로저를 작성
//				wardDone = make(chan interface{}) // 피후견인 고루틴으로 전달할 중단신호용 채널 생성
//				// 모니터링할 고루틴 시작, 스튜어드가 멈추거나 와드 고루틴을 멈추게 하고자 하는 경우
//				// 와드 고루틴이 멈추기 원하기 때문에 논리적 or로 두 가지 done 채널을 모두 감싸준다.
//				//wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2) // 이거 확인해봐야됨; orChan인가;;
//				wardHeartbeat = startGoroutine(common.OrDone(wardDone, done), timeout/2) // 이거 확인해봐야됨; orChan인가;;
//
//			}
//			startWard()
//			pulse := time.Tick(pulseInterval)
//
//		monitorLoop:
//			for {
//				timeoutSignal := time.After(timeout)
//
//				for { // 스튜어드가 자체적으로 내부 루프를 보낼 수 있도록 하는 내부 루프
//					select {
//					case <-pulse:
//						select {
//						case heartbeat <- struct{}{}:
//						default:
//						}
//					case <-wardHeartbeat: // 와드의 하트비트를 받으면 계속 모니터링 루프를 진행
//						continue monitorLoop
//					case <-timeoutSignal: // 타임아웃시 와드 중단시키고 새로운 와드 고루틴이 시작되도록 요청
//						log.Println("steward: ward unhealthy; restarting")
//						close(wardDone)
//						startWard()
//						continue monitorLoop
//					case <-done:
//						return
//					}
//				}
//			}
//		}()
//
//		return heartbeat
//	}
//}
//
//// for 루프가 좀 바쁘긴 하지만...
