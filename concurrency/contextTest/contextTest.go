package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

//func main() {
//	var wg sync.WaitGroup
//	done := make(chan interface{})
//	defer close(done)
//
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		if err := printGreeting(done); err != nil {
//			fmt.Printf("%v", err)
//			return
//		}
//	}()
//
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		if err := printFarewell(done); err != nil {
//			fmt.Printf("%v", err)
//			return
//		}
//	}()
//
//	wg.Wait()
//}
//
//func printGreeting(done <-chan interface{}) error {
//	greeting, err := genGreeting(done)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("%s world!\n", greeting)
//
//	return nil
//}
//
//func printFarewell(done <-chan interface{}) error {
//	farewell, err := genFarewell(done)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("%s world!\n", farewell)
//	return nil
//}
//
//func genGreeting(done <-chan interface{}) (string, error) {
//	switch locale, err := locale(done); {
//	case err != nil:
//		return "", err
//	case locale == "EN/US":
//		return "hello", nil
//	}
//	return "", fmt.Errorf("unsupported locale")
//}
//
//func genFarewell(done <-chan interface{}) (string, error) {
//	switch locale, err := locale(done); {
//	case err != nil:
//		return "", err
//	case locale == "EN/US":
//		return "goodbye", nil
//	}
//	return "", fmt.Errorf("unsupported locale")
//}
//
//func locale(done <-chan interface{}) (string, error) {
//	select {
//	case <-done:
//		return "", fmt.Errorf("canceled")
//	case <-time.After(1 * time.Minute):
//	}
//	return "EN/US", nil
//}

// 아래 코드는 ㅈ
func main() {
	// 2번 예제
	ProcessRequest("Jane", "auth123")

	// 1번 예제
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printGreeting(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v \n", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()

	wg.Wait()
}

func printGreeting(ctx context.Context) error {
	greeting, err := genGreeting(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(ctx context.Context) error {
	farewell, err := genFarewell(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genGreeting(ctx context.Context) (string, error) {
	// context will be canceled after 1 Second. d
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(ctx context.Context) (string, error) {
	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

// check if deadline is given (실행 결과 자체는 동일, 그러나 빠르게 실패할 수 있음)
func locale(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}

// 유일한 문제는 하위의 호출 그래프가 얼마나 오래 걸리는지 알고 있어야 한다는 점...
// 이로 인해 context 패키지가 사용하는 Context용 데이터 저장소를 사용하게 된다. 함수가 고루틴과 Context를 생성할 때, 많은 경우에
// 요청을 처리할 프로세스를 시작하며, 스택의 아래쪽에 있는 함수는 요청에 댛나 정보를 필요로 한다는 것을 잊지 말자. 다음은 Context에
// 데이터를 저장하고 조회하는 예제이다.
// 추가적으로 컨텍스트의 키 타입을 재정의해 외부에서 안전하게 만드는 것 까지..

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(c context.Context) string {
	return c.Value(ctxUserID).(string)
}

func AuthToken(c context.Context) string {
	return c.Value(ctxAuthToken).(string)
}

func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"hendling response for %v (%v) \n",
		UserID(ctx),
		AuthToken(ctx),
	)
}

// 이 기법에는 문제가 있다. HandleResponse가 response라는 다른 패키지에 있다고 할 떄, 그리고 ProcessRequest는 process 패키지 내에
// 있다고 할 때, Handle Response를 호출하기 위해 response 패키지를 가져와야 하지만, HandelResponse는 process 패키지에 정의된
// 접근자 함수에 접근할 수 없다. process를 입포트 할 때 순환 의존성을 발생시키기 때문이다. 기본적으로 Context에 키를 저장하는데 사용되는
// 타입이 비공개이기 때문에!!!!
// context 패키지는 다소 논란의 여지가 있다. 임의 데이터를 저장할 수 있는 기능과 타입에 안전하지 않는 방식의 데이터 저장은 논란이 되어왔다.
// 더 큰 문제는 개발자가 Context의 인스턴스에 저장해야만 하는 특성이다. 어떤식으로 데이터를 저장해야 하는가..?
// 저자는 다음과 같은 지침을 권장한다.
// 1. 데이터가 API나 프로세스 경계를 통과해야 한다.
// 2. 데이터는 immutable 해야 한다.
// 3. 데이터는 단순한 타입으로 변해야 한다.
// 4. 데이터는 메서드를 포함하지 않는다.
// 5. 데이터는 주도적이지 않아야 한다. 그런 것은 매개변수 영역으로 넘어가야 한다.
// 이는 모두 경험적인 규칙이다. 이 모든 것을 위반한다면,수행하는 작업을 다시 살펴보자...
// 또한 이 데이터가 거쳐가야 할 레이어의 수도 고려해야 한다.
