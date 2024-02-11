package algorithms

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSuperset(t *testing.T) {
	input := []int{1, 2, 3}
	output := [][]int{
		{1},
		{2},
		{3},
		{1, 2},
		{1, 3},
		{2, 3},
		{1, 2, 3},
	}
	assert.Equal(t, GetSuperset(input), output, "Wrong")
}
