package main

import "fmt"

// 때로는 시스템에서 서로 다른 부분의 채널들로 작업하게 되는 경우가 있다.
// 작업 중인 코드가 done 채널을 통해 취소될 때, 채널이 어떻게 동작할지 단언할 수는 없다.
// 고루틴이 취소됐다는 것이, 읽어오는 채널 역시 취소되었음을 의미하지는 않는다.
// 고루틴 누수에서 했던 것 같이, 채널에서 읽어오는 부분을 done 채널을 select 구문으로 감싸야 한다.
// 다음의 코드가 필요하다.
func main() {
	done := make(chan interface{})
	myChan := make(chan interface{})
	// 다음 코드를 더 간단하게 바꿔봅시다.
	// loop:
	//for {
	//	select {
	//	case <-done:
	//		break loop
	//	case maybeVal, ok := <-myChan:
	//		if ok == false {
	//			return
	//		}
	//		// val로 무언가를 한다
	//	}
	//}
	for val := range orDone(done, myChan) {
		// val로 무언가를 한다
		fmt.Print(val)
	}
}

// orDone - 캡슐화 버전
func orDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}
