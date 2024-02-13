package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var n int
var nums []int
var check []int // 순열 체크용
var answer int = -99999

func Dfs(perm []int, lv int) {
	if lv == n {
		fmt.Printf("perm: %v\n", perm)
		return
	}

	for i := 0; i < n; i++ {
		if check[i] == 1 {
			continue
		}
		check[i] = 1
		perm = append(perm, nums[i])
		Dfs(perm, lv+1) // 마킹한 미래를 갖다 쓴다
		check[i] = 0    // 이 미래를 취소하고 다음 미래를 모색한다
		lastIdx := len(perm) - 1
		perm = perm[0:lastIdx]
	}
}

func GetTotal(nums []int) int {
	tot := 0
	for i := 1; i < n; i++ {
		tot += Abs(nums[i-1] - nums[i])
	}
	return tot
}

func Abs(num int) int {
	if num >= 0 {
		return num
	} else {
		return -num
	}
}
func main() {
	fmt.Scanln(&n)

	nums = make([]int, n)
	check = make([]int, n)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	for i := 0; i < n; i++ {
		scanner.Scan()
		nums[i], _ = strconv.Atoi(scanner.Text())
	}

	var emptySlice []int
	Dfs(emptySlice, 0)
	fmt.Println(answer)
}
