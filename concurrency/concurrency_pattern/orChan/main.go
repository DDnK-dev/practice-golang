package main

import (
	"fmt"
	"time"
)

func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} { // 가변 slice channel 받아 channel 리턴
		switch len(channels) {
		case 0: // 재귀함수이므로... 종료 기준
			return nil
		case 1: // 가변슬라이스에 요소 하나만 있으면 하나 리턴
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() { // 함수 핵심부분, 재귀가 발생하는 부분. 채널들에서 차단없이 메시지 대기하도록 고루틴 생성
			defer close(orDone)

			switch len(channels) {
			case 2: // 재귀 방식을 사용하여 or에 대한 모든 재귀호출은 최소 두개 채널을 가짐. 최적화용
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			// 슬라이스 세번째 인덱스 이후 위치한 모든 채널에서 재귀적으로 or 채널을 만든 다음 이 중에서 select 수행
			// 첫 번째 신호가 리턴되는 것으로부터 트리를 형성하기 위해 나머지 슬라이스를 or 채널로 분해
			// 또한 고루틴들이 트리를 위쪽과 아래쪽 모두 빠져나올 수 있도록 orDone 채널을 전달
			// 즉 여러 개 채널을 한개의 채널로 결합해, 여러 채널 중 하나라도 닫히거나 데이터가 쓰여지면
			// 모든 채널이 닫히도록 할 수 있는 매우 간결한 함수....
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}

	// 설정된 지속 시간 후에 닫히는 여러 채널들을 받아서, 이들을 or 함수를 사용해 단일 채널로 결합하고 닫는 예제
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
