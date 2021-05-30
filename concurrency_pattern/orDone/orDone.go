package main

import "fmt"

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
