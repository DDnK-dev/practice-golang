/*
buffered-channel 이해를위한 적절한 예시
*/

package main

import "fmt"

// bufffered channel 예시
//func main() {
//	var stdoutBuff bytes.Buffer
//	defer stdoutBuff.WriteTo(os.Stdout)
//
//	intStream := make(chan int, 4)
//	go func() {
//		defer close(intStream)
//		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
//		for i := 0; i < 4; i++ {
//			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
//			intStream <- i
//		}
//	}()
//
//	for integer := range intStream {
//		fmt.Fprintf(&stdoutBuff, "Received %v. \n", integer)
//	}
//}

// 채널 소유 및 소비 예시
func main() {
	chanOwner := func() <-chan int{
	resultStream := make(chan int, 5)
	go func() {
		defer close(resultStream)
		for i := 0; i <= 5; i++ {
			resultStream <- i
		}
	}()
	return resultStream
	}

	resultStream := chanOwner()
	for result := range resultStream {
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Printf("Done receiving!")
}
