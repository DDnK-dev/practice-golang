package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 티케팅을 생각해보자. 빠르게 요청을 받아오기 위하여 요청을 여러개의 핸들러에 뿌려 다량의 결과를 요청하고, 그 중 가장 빠른 것 하나만을
// 취할 수 있을 것이다. 다음 예제는 10개의 핸들러로 요청을 복제한다.
func main() {
	exampleSimulater()
}

func exampleSimulater() {
	done := make(chan interface{})
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go copyRequestToTenHandler(done, i, &wg, result)
	}

	firstReturned := <-result
	close(done)
	wg.Wait()

	fmt.Printf("Recieved and answer from #%v\n", firstReturned)
}
func copyRequestToTenHandler(
	done <-chan interface{},
	id int,
	wg *sync.WaitGroup,
	result chan<- int,
) {
	started := time.Now()
	defer wg.Done()

	// random work Load simulation
	simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second
	select {
	case <-done:
	case <-time.After(simulatedLoadTime):
	}

	select {
	case <-done:
	case result <- id:
	}

	took := time.Since(started)
	// show how long does the handler takes
	if took < simulatedLoadTime {
		took = simulatedLoadTime
	}
	fmt.Printf("%v took %v\n", id, took)
}
