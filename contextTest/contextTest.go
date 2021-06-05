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

// 이 기법에는 문제가 있다. HandleResponse가 response라는 다른 패ㅋ지에 있다고 할 떄,.
