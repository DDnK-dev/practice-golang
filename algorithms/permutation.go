package algorithms

import (
	"math"
	"strconv"
)

func solution(numbers string) int {
	check := make([]bool, len(numbers))
	var buff []rune

	alreadySeen := make(map[int]struct{})
	ans := 0
	f := func(dfs interface{}, cnt int, lv int) { // dfs
		// 탈출조건
		if cnt >= lv { // 버퍼가 다 차면
			intVal, _ := strconv.Atoi(string(buff))
			if _, ok := alreadySeen[intVal]; ok { // 기존에 검사한 값이면 패스
				return
			}
			alreadySeen[intVal] = struct{}{}
			if checkPermutation(intVal) {
				ans += 1
			}
			return
		}
		for i := range numbers {
			if check[i] {
				continue
			}
			buff = append(buff, rune(numbers[i]))
			check[i] = true
			dfs.(func(interface{}, int, int))(dfs, cnt+1, lv)
			check[i] = false
			buff = buff[:len(buff)-1]
		}
	}

	for i := range numbers {
		buff = make([]rune, 0, i+1)
		f(f, 0, i+1)
	}

	// permutation list
	return ans
}

func checkPermutation(val int) bool {
	if val < 2 {
		return false
	} else if val == 2 {
		return true
	}

	upper := int(math.Ceil(math.Sqrt(float64(val))))
	for i := 2; i <= upper; i++ {
		if val%i == 0 {
			return false
		}
	}
	return true
}
