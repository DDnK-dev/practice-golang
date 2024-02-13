package fullsearch

import "math"

// 프로그래머스 : https://school.programmers.co.kr/learn/courses/30/lessons/42842
// 노란색 사이즈가 x, y 일때,
// yellow = xy
// brown = 2y + 2x + 4
// 이떄, 가로 길이(x)는 세로 길이와 같거나, 길다

func solution1(brown int, yellow int) []int {
	// 노란색 사이즈가 x, y 일때,
	// yellow = xy
	// brown = 2y + 2x + 4
	// 이떄, 가로 길이(x)는 세로 길이와 같거나, 길다
	x := int(math.Floor(math.Sqrt(float64(yellow)))) // 엄밀하지 않으므로 정렬 필요
	res := make([]int, 2)
	for i := x; x <= yellow; i++ {
		if yellow%i == 0 {
			y := yellow / i
			if 2*i+2*y+4 == brown {
				res[0] = y + 2
				res[1] = i + 2
				break
			}
		}
	}
	if res[0] < res[1] {
		return []int{res[1], res[0]}
	}
	return res
}
