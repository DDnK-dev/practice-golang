package dfs

import (
	"fmt"
	"sort"
)

func solutionTicket(tickets [][]string) []string {
	// 일단 map을 만들자.. (for문 계속 돌아야 할 것 같음)
	tMap := make(map[string][]string)
	for _, tt := range tickets {
		t := tt
		if _, ok := tMap[t[0]]; !ok {
			tMap[t[0]] = make([]string, 0, 4) // magic number
		}
		id := fmt.Sprintf("%s%d", t[1], len(tMap[t[0]]))
		tMap[t[0]] = append(tMap[t[0]], id)
	}
	var dfs func(string, int)
	ress := [][]string{}
	res := make([]string, 0, len(tickets)+1)
	res = append(res, "ICN")
	seen := make(map[string]struct{})

	dfs = func(src string, cnt int) {
		if cnt >= len(tickets) {
			newRes := make([]string, len(tickets)+1)
			copy(newRes, res)
			ress = append(ress, newRes)
			return
		}
		var dests []string
		var ok bool
		if dests, ok = tMap[src]; !ok {
			return // 잘못왔네? // 문제 전제상 이런 경우 없어야 정상
		}
		dests = tMap[src]
		for i := range dests {
			if _, ok := seen[src+dests[i]]; ok {
				continue
			}
			seen[src+dests[i]] = struct{}{}
			res = append(res, dests[i][:3])
			dfs(dests[i][:3], cnt+1)
			res = res[:len(res)-1]
			if _, ok := seen[src+dests[i]]; ok {
				delete(seen, src+dests[i])
			}
		}
	}
	dfs("ICN", 0)
	if len(ress) > 1 {
		sort.Slice(ress, func(a, b int) bool {
			for i := range ress[0] {
				if ress[a][i] < ress[b][i] {
					return true
				} else if ress[a][i] > ress[b][i] {
					return false
				}
			}
			return true
		})
	}
	return ress[0]
}
