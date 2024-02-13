package dfs

// 문제: https://school.programmers.co.kr/learn/courses/30/lessons/43165
// 더하거나 뺴서 target을 만들 수 있는 경우의 수는?

func solution(numbers []int, target int) int {
	// dfs로 풀어보자 (모든 경우의 수르 구해야 하므로)
	// 여기서만 쓸 함수니까 dfs
	var f func(int, int) int
	f = func(cnt int, sum int) int {
		if cnt == len(numbers) {
			if sum == target {
				return 1
			}
			return 0
		}
		return f(cnt+1, sum+numbers[cnt]) + f(cnt+1, sum-numbers[cnt])
	}
	return f(0, 0)
}
