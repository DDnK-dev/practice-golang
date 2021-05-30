/*
This package shows how to use channel to terminate child goroutine
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

//func main() {
//	doWork := func(
//		done <-chan interface{},
//		strings <-chan string,
//	) <-chan interface{} {
//		terminated := make(chan interface{})
//		go func() {
//			defer fmt.Println("doWork exited.")
//			defer close(terminated)
//			for {
//				select {
//				case s := <-strings:
//					fmt.Println(s)
//				case <-done:
//					return
//				}
//			}
//		}()
//		return terminated
//	}
//
//	done := make(chan interface{})
//	terminated := doWork(done, nil) // nil channel 을 인자로 줌 -> 가본적으로 done 없으면 무한대기
//
//	go func() {
//		// 1초 후에 작업 취소
//		time.Sleep(time.Second * 1)
//		fmt.Println("Canceling doWork goroutine...")
//		close(done)
//	}()
//
//	<-terminated // doWork 에서 생성된 고루틴과 main 고루틴을 조인
//	fmt.Println("Done")
//}

// 채널에 값으 쓰려는 시도를 차단하는 고루틴의 경우는 어떻게 하나??
//func main() {
//	newRandStream := func() <-chan int {
//		randStream := make(chan int)
//		go func() {
//			defer fmt.Println("newRandStream closure exited. ")
//			defer close(randStream)
//			for {
//				randStream <- rand.Int()
//			}
//		}()
//		return randStream
//	}
//
//	randStream := newRandStream()
//	fmt.Println("3 random ints:")
//	for i := 1; i <= 3; i++ {
//		fmt.Printf("%d: %d\n", i, <-randStream)
//	}
//}
// 출력 결과로부터 defer fmt.Println("newRandStream closure exited. ") 코드 실행되지 않음을 확인
// 따라서 해결책은 생산자 고루틴에게 종료를 알리는 채널을 제공하는 것

func main() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)

	time.Sleep(1 * time.Second)
}
