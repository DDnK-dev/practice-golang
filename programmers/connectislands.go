package programmers

// https://programmers.co.kr/learn/courses/30/lessons/42861?language=go
// 최소의 비용으로 모든 섬이 서로 통행 가능

type edge struct {
	dest int
	cost int
}

func SolutionIslands(n int, costs [][]int) int {
	// 일단 map 만들자
	var ret int
	visitedSet := make(map[int]struct{})
	islandMap := make(map[int][]edge)
	for i, _ := range costs {
		newEdge := edge{dest: costs[i][1], cost: costs[i][2]}
		islandMap[costs[i][0]] = append(islandMap[costs[i][0]], newEdge)
	}

	// 0부터 시작
	visitedSet[0] = struct{}{}
	for cnt := 0; cnt < n; cnt = cnt + 1 {
		target := edge{}
		minCost := int(^uint(0) >> 1) // set int.Inf
		for k, _ := range visitedSet {
			for i := range islandMap[k] {
				e := islandMap[k][i]
				_, exist := visitedSet[e.dest]
				if e.cost < minCost && !exist {
					target = e
				}
			}
		}
		ret += target.cost
		visitedSet[target.dest] = struct{}{}
	}

	return ret
}

//0 1 1
//0 2 2
//1 2 5
//1 3 1
//2 3 8
