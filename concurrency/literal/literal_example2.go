package literal

import (
	"fmt"
	"os"
)

type Writer func(string)

func writerHello(writer Writer) {
	writer("Hello World")
}

func main() {
	f, err := os.Create("test.txt")
	if err != nil {
		fmt.Println("Failed to create a file")
		return
	}
	defer f.Close()

	writerHello(func(msg string) {
		fmt.Fprintf(f, msg)
	})
}
