/*
 해당 코드는 for loop와 goroutine의 스케줄링 관계를 파악하는데 도움을 준다.
예상하는 출력결과와 실제 출력결과가 다르다
*/
package main

import (
	"fmt"
	"sync"
)

//func main() {
//	var wg sync.WaitGroup
//	for _, salutation := range []string{"hello", "greetings", "good day"} {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			fmt.Println(salutation)
//		}()
//	}
//	wg.Wait()
//}

func main() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()
}
