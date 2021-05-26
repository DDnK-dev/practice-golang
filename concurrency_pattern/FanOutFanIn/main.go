package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	FindPrimeNum()
}

// FindPrimeNum finds prime number within given stream. basically, it's very slow algorithm
func FindPrimeNum() {
	rand := func () interface {} {return rand.Intn(500000000)}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))
}

// primeFinder finds prime number, dividing values with smaller numbers
func primeFinder(
	done <-chan interface{},
	valueStream <-chan int,
	)<- chan int{
	primeStream := make(chan int)
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
func repeatFn (
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
	valueStream <-chan int,
	num int,
) <-chan int {
	takeStream := make(chan int)
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
