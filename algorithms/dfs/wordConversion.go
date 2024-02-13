package dfs

// 문제: https://school.programmers.co.kr/learn/courses/30/lessons/43163?language=go
// 단어 변환
// 단어 begin에서 target으로 변환하는 가장 짧은 변환 과정의 길이를 반환하라
// 단, 한 번에 한 개의 알파벳만 바꿀 수 있고, words에 있는 단어로만 변환할 수 있다.
// words에 target이 없으면 0을 반환한다.

type cache struct {
	c map[string]map[string]uint8
}

func newCache() *cache {
	return &cache{
		c: make(map[string]map[string]uint8),
	}
}

func (c *cache) Set(a, b string, cost int) {
	if a < b {
		a, b = b, a
	}
	if _, ok := c.c[a]; !ok {
		c.c[a] = make(map[string]uint8)
	}
	c.c[a][b] = uint8(cost)
}

func (c *cache) Get(a, b string) int {
	if a < b {
		b, a = a, b
	}
	if _, ok := c.c[a]; !ok {
		return c.get(a, b)
	} else if _, ok := c.c[a][b]; !ok {
		return c.get(a, b)
	}
	return int(c.c[a][b])
}

// 내부적으로 데이터가 없을때 값 세팅하는 함수
func (c *cache) get(a, b string) int {
	dist := getHammingDistance(a, b)
	c.Set(a, b, dist)
	return dist
}

func solution2(begin string, target string, words []string) int {
	minCost := int(^uint(0) >> 1) // inf

	// 성능 향상을 위해 노드간 hamming distance를 기억하는 맵을 하나 만든다
	memCache := newCache()

	// 완전 탐색을 위한 dfs 선언
	// 루프 방지를 위해 이미 간 곳은 가지 않는다.
	var f func(string, int)
	seen := make([]bool, len(words))
	f = func(src string, cnt int) {
		if cnt > len(words) { // 지정 횟수 초과
			return
		}
		if src == target && minCost > cnt {
			minCost = cnt
		}
		for i := range words {
			if seen[i] {
				continue
			}
			seen[i] = true
			dist := memCache.Get(src, words[i])
			if dist == 1 {
				f(words[i], cnt+1)
			}
			cnt -= 1
			seen[i] = false
		}
	}

	f(begin, 0)
	return minCost
}

// cache 헬퍼 함수
func getHammingDistance(a, b string) int {
	distance := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			distance += 1
		}
	}
	return distance
}
