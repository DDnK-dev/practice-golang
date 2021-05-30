/*
bridge 채널.
연속된 채널로 부터 값을 사용하고 싶을 때가 있다.
채널을 하나로 병합하는것과는 달리, 서로 다른 출처에서부터 순서대로 값을 쓴다는 점을 유의하자.
예를 들어, 제한 패턴을 따르고, 채널에 대한 소유권을 해당 채널에 쓰는 고루틴으로 이전한다면,
매번 새로운 채널이 생성되게 될 것이다.
이 점은 실질적으로 연속된 채널을 가진다는 것을 의미한다.

  이번 예제에서는 bridge 채널을 만들고, 각 채널들에서 오는 값을 하나의 채널에 송신한다.
이를 위해 하나의 요소가 쓰여진 연속된 10개의 채널을 생성하고, 이 채널들을 bridge 함수로 전달한다.
*/

package main

import (
	"fmt"
	"practice-golang/common"
)

func main() {
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream) // 닫힌 채널에 값을 넣을수는 없지만, 뺄 수는 있다는 것을 유의하자
				chanStream <- stream
			}
		}()
		return chanStream
	}
	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v", v)
	}
}

func bridge(
	done <-chan interface{},
	chanStream <-chan <-chan interface{},
) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range common.OrDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}
