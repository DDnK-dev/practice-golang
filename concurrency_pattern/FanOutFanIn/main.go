package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	rand := func() interface{} { return rand.Intn(5000000) }

	FindPrimeNum(rand)
	FindPrimeNumFanOut(rand)
}

// FindPrimeNum finds prime number within given stream. basically, it's very slow algorithm
func FindPrimeNum(rand func() interface{}) {
	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v \n", time.Since(start))
}

// FindPrimeNumFanOut is fan out version FindPrimeNum function
func FindPrimeNumFanOut(rand func() interface{}) {
	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Primes")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v \n", time.Since(start))
}

// primeFinder finds prime number, dividing values with smaller numbers
func primeFinder(
	done <-chan interface{},
	valueStream <-chan int,
) <-chan interface{} {
	primeStream := make(chan interface{})
	go func() {
		for {
			select {
			case <-done:
				return
			case value := <-valueStream:
				// Logic
				isPrime := true
				for i := 2; i < value; i++ {
					if value%i == 0 {
						isPrime = false
						break
					}
				}
				if isPrime {
					primeStream <- value
				}
			}
		}
	}()
	return primeStream
}

// repeatFn generates stream repeatedly using given function
func repeatFn(
	done <-chan interface{},
	fn func() interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

// toInt is assertion logic. converts type of elements in stream into Int type
func toInt(
	done <-chan interface{},
	valueStream <-chan interface{},
) <-chan int {
	stringStream := make(chan int)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(int):
			}
		}
	}()
	return stringStream
}

// take reads value from channel num times
func take(
	done <-chan interface{},
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

// fainIn
func fanIn(
	done <-chan interface{},
	channels ...<-chan interface{},
) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}
	wg.Add(len(channels))

	for _, c := range channels {
		go multiplex(c)
	}

	//wait until all reads end
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}
