package main

import (
	"context"
	"log"
	"os"
	"sync"
)

// 토큰을 수용할수 있는 버퍼를 두고, 일정 텀으로 토큰을 버퍼에서 제거해 특정 수만큼만 요청을 받아들이는 로직
// API에 대한 요청을 받는다고 가정하고, 다음의 클라이언트를 사용한다고 가정한다.

func Open() *APIConnection {
	return &APIConnection{}
}

type APIConnection struct{}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	// 여기서 무언가 작업이 이루어진다
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	// 여기서 무언가 작업이 이루어진다.
	return nil
}

// 이제 이 API에 접근할 수 있는 간단한 드라이버를 만들 것이다. 10개의 파일을 읽고 10개의 주소를 확인해야 한다.
// 서로 관련이 없으므로 드라이버는 API 호출을 동시에 수행할 수 있다.
func main() {
	case1()
}

// 모든 요청이 거의 동시에 처리됨을 확인할 수 있음. 속도제한이 없는 상태에서는 무한루프가 발생했을 떄 어마어마한 청구서를 받을 수 있다...
func case1() {
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}

	wg.Wait()
}
