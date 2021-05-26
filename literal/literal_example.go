package literal

import "fmt"

func CaptureLoop () {
	f := make([]func(), 3) 		//len 3의 func slice 생성
	fmt.Println("ValueLoop")
	for i:= 0; i < 3; i++ { // 각 요소에 함수를 저장
		f[i] = func() {
			fmt.Println(i)
		}
	}

	for i := 0; i < 3; i++{ // 각 요소에 저장된 함수 실행; Capturing i
		f[i]()
	}
}

func CaptureLopop2() {
	f := make([]func(), 3)
	fmt.Println("ValueLoop2")
	for i := 0; i < 3; i++ {
		v := i	// copy i's value to v
		f[i] = func() {
			fmt.Println(v)
		}
	}
	for i := 0; i < 3; i++ {
		f[i]()
	}
}

func main() {
	CaptureLoop()
	CaptureLopop2()
}
