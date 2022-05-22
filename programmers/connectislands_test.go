package programmers_test

import (
	"programmers"
	"testing"
)

func TestIslands(t *testing.T) {
	n := 4
	costs := [][]int{
		{0, 1, 1}, {0, 2, 2}, {1, 2, 5}, {1, 3, 1}, {2, 3, 8},
	}
	s := programmers.SolutionIslands(n, costs)
	if s != 4 {
		t.Errorf("failed to process solution")
	}
}
