/*
리눅스에서 tee는 std input에서 읽어 std output과 파일에 쓰는 명령어.
예를들어, echo "hello" | tee OUTFILE

tee channel은 여기서 이름을 따온 것. 채널에서 들어오는 값을 분리해 별개의 두 영역으로 보내고자 할 때 사용.
사용자 명령 채널을 예로 들면, 채널에서 사용자 명령 스트림을 가져와 실행자에게보내고, 감사 프로세스에 이걸 보내고...

읽어들인 채널을 전달할ㄹ 수 있고 ,동일한 값을 얻어오는 별개의 채널 두 개를 리턴한다.
*/
package main

import (
	"fmt"
	"practice-golang/common"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, common.Take(done, common.Repeat(done, 1, 2), 4))

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v2\n", val1, <-out2)
	}
}

func tee(
	done <-chan interface{},
	in <-chan interface{},
) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range orDone(done, in) {
			var out1, out2 = out1, out2 // 지역변수를 사용해 밖에있는거 가리기
			for i := 0; i < 2; i++ {    // out1, out2가 서로를 가리지 않도록 두 번
				select {
				case <-done:
				case out1 <- val:
					out1 = nil // 채널에 쓴 이후 로컬 복사본을 nil로 설정해 추가적인 기록이 안 되도록 차단
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
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
