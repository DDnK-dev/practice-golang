package main

import (
	"fmt"
	"math/rand"
)

func main() {
	//DiscreteSetToStream()
	//DiscreteSetToFnStream()
	DiscreteSetToStreamWithAssertion()
}

// DiscreteSetToStream shows how to convert discrete set to stream using channel generator
func DiscreteSetToStream() {
	// repeat continues sending values to channel until get signal done
	repeat := func(
		done <-chan interface{},
		values ...interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}() // parameter null
		return valueStream
	}

	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}
}

// DiscreteSetToFnStream call function repeatedly
func DiscreteSetToFnStream() {
	repeatFn := func(
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

	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} { return rand.Int() }

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}
}

// DiscreteSetToStreamWithAssertion shows how to apply type assertion.
// In this case, it use ToString function to show type assertion
func DiscreteSetToStreamWithAssertion() {
	repeat := func(
		done <-chan interface{},
		values ...interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}() // parameter null
		return valueStream
	}

	toString := func(
		done <-chan interface{},
		valueStream <-chan interface{},
	) <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case stringStream <- v.(string):
				}
			}
		}()
		return stringStream
	}
	done := make(chan interface{})
	defer close(done)

	var message string
	for token := range toString(done, take(done, repeat(done, "I", "am"), 5)) {
		message += token
	}
	fmt.Printf("message: %s...", message)
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

func repeat(
	done <-chan interface{},
	values ...interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}() // parameter null
	return valueStream
}

func toString(
	done <-chan interface{},
	valueStream <-chan interface{},
) <-chan string {
	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}
	}()
	return stringStream
}
