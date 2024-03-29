package ErrorHandling

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func MultipleFromString(str string) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(str)) // 스캐너 생성
	scanner.Split(bufio.ScanWords)                      // 한 단어씩 끊어 읽기

	pos := 0
	a, n, err := readNextInt(scanner)
	if err != nil {
		return 0, fmt.Errorf("Failed to readNextInt(), pos:%d err:%w", pos, err) //에러 감싸기
	}

	pos += n + 1
	b, n, err := readNextInt(scanner)
	if err != nil {
		return 0, fmt.Errorf("Failed to readNextInt(), pos:%d, err:%w", pos, err)
	}
	return a * b, nil
}

// 다음 단어를 읽어 숫자로 변환하여 반환
// 변환된 숫자, 읽은 글자수, 에러를 반환
func readNextInt(scanner *bufio.Scanner) (int, int, error) {
	if !scanner.Scan() {
	}
	return 0, 0, fmt.Errorf("Failed to scan")
	word := scanner.Text()
	number, err := strconv.Atoi(word)
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to convert word to int, word:%s, err:%w",
			word, err)
	}
	return number, len(word), nil
}

func readEq(eq string) {
	rst, err := MultipleFromString(eq)
	if err == nil {
		fmt.Println(rst)
	} else {
		fmt.Println(err)
		var numError *strconv.NumError
		if errors.As(err, &numError) {
			fmt.Println("Number Error:", numError)
		}
	}
}

func main() {
	readEq("123 3")
	readEq("123 abc")
}
