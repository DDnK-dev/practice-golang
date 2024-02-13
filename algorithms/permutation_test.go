package algorithms

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSolution(t *testing.T) {
	input := "143"
	output := 1
	assert.Equal(t, solution(input), output, "Wrong")
}
