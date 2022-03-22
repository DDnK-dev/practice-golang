package ChannelAndContext

import (
	"fmt"
	"sync"
	"time"
)

type Car struct {
	Body  string
	Tire  string
	Color string
}

var wg3 sync.WaitGroup
var startTime = time.Now()

func main() {
	tireCh := make(chan *Car)
	paintCh := make(chan *Car)

	fmt.Printf("Start Factory\n")

	wg3.Add(3)
	go MakeBody(tireCh)
	go InstallTire(tireCh, paintCh)
	go PaintCar(paintCh)

	wg3.Wait()
	fmt.Println("Close the factory")
}

func MakeBody(tireCh chan *Car) {
	tick := time.Tick(time.Second)
	after := time.After(10 * time.Second)
	for {
		select {
		case <-tick:
			// Make a body
			car := &Car{}
			car.Body = "Sports car"
			tireCh <- car
		case <-after:
			close(tireCh)
			wg3.Done()
			return
		}
	}
}

func InstallTire(tireCh, paintCh chan *Car) {
	for car := range tireCh {
		//Make a body
		time.Sleep(time.Second)
		car.Tire = "Winter tire"
		paintCh <- car
	}
	wg3.Done()
	close(paintCh)
}

func PaintCar(paintCh chan *Car) {
	for car := range paintCh {
		// Make a body
		time.Sleep(time.Second)
		car.Color = "Red"
		duration := time.Now().Sub(startTime) // 경과 시간 출력
		fmt.Printf("%.2f Complete Car: %s %s %s \n", duration.Seconds(),
			car.Body, car.Tire, car.Color)
	}
	wg3.Done()
}

//...........z
