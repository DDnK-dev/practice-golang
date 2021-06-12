package main

import (
	"context"
	"golang.org/x/time/rate"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

// rate 패키지는 time.Duration을 Limit로 변환하는 것을 돕기 위해 도우미 메서드인 Every를 정의한다.
// 의미가 있지만... 여기선 속도 제한을 요청 간격이 아니라 측정시간 당 작업 수라는 면에서 논의하고자 한다.  이는
// rate.Limit(events/timePeriod.Second()) 로 표현이 가능하다.
// 매번 이를 입력하고 싶지 않으므로... 또 interval이 0이면 rate.Inf로 리턴하는 특수논리를 지니고 있는걸 고려해 이를 구현한다.

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

// 일자별로 Limit를 다르게 두고 싶다고 해보자. 이 때는 여러개의 속도 제한기를 만들고, 이를 하나로 결합해 관리하는것이 나을 것이다.
// 이 예제는 multiLimiter라는 간단한 통합된 속도 제한기를 보여준다.

type RateLimiter interface {
	Wait(ctx context.Context) error
	Limit() rate.Limit
}

func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

type multiLimiter struct {
	limiters []RateLimiter
}

// Wait 메서드는 루프를 통해 모든 자식 속도 제한기를 순회하며 각 자식 속도 제한기를 호출. 대기일 수도 아닐수도 있다.
// 요청의 각 속도제한기에 통지해야만 토큰 버킷을 줄일 수 있음. 결과적으로 가장 긴 대기 시간동안 기다리게 된다.
func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

// 이를 이용하는 예제

func Open() *APIConnection {
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)   // 대량 요청을 처리하지 않는 제한
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10) // 초기 풀을 제공하기 위해 10개 대량 요청처리하는 제한
	return &APIConnection{
		rateLimiter: MultiLimiter(secondLimit, minuteLimit), // 두 제한을 결합해 API Connection의 주 속도제한기로 설정
	}
}

type APIConnection struct {
	rateLimiter RateLimiter
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil { //
		return err
	}
	// 여기서 무언가 작업이 이루어진다
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// 여기서 무언가 작업이 이루어진다.
	return nil
}

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

// 결과는 다음과 같다. 11번째 요청까지 초당 두 건을 요청.
// 이 시점부터 6초마다 요청하기 시작.
// 11번째 요청은 2초 후에 발생....
// 시스템 속도의 관점에서 생각해보자!!!!!!!!
