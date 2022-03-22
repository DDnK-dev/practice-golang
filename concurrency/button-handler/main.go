/*
GUI 프로그래밍중이라 가정하고, 마우스 클릭시 func을 실행하는 핸들러를 만든다.
이는 BroadCasting을 확인하기 위함이다.
*/

package main

import (
	"fmt"
	"sync"
)

type Button struct {
	Clicked *sync.Cond
}

func main() {
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	//편의 함수 설정. 조건의 신호를 처리.
	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	// 헨들러 등록
	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse Clicked")
		clickRegistered.Done()
	})

	//핸들러 실행
	button.Clicked.Broadcast()

	clickRegistered.Wait()
}
