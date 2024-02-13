package dfs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSolution(t *testing.T) {
	assert.Equal(t, solution2("hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"}), 4, "Wrong")
	assert.Equal(t, solution2("hit", "cog", []string{"hot", "dot", "dog", "lot", "log"}), 0, "Wrong")
}
