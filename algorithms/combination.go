package algorithms

// GetSuperset 공집합을 제외한 모든 부분집합을 가져온다
func GetSuperset(numbers []int) [][]int {
	ress := make([][]int, 0, 1024)
	superset(numbers, make([]bool, len(numbers)), 0, ress)
	return ress
}

func superset(numbers []int, sel []bool, idx int, ress [][]int) {
	if idx >= len(numbers) {
		res := make([]int, 0, len(numbers))
		for i := range sel {
			if sel[i] {
				res = append(res, numbers[i])
			}
		}
		ress = append(ress, res)
		return
	}
	newSel := make([]bool, len(sel))
	copy(newSel, sel)
	newSel[idx] = true
	superset(numbers, newSel, idx+1, ress)
	newSel[idx] = false
	superset(numbers, newSel, idx+1, ress)
}
