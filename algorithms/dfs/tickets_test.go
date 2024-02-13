package dfs

import (
	"testing"
)

func TestSolutionT(t *testing.T) {
	output := solutionTicket([][]string{{"ICN", "JFK"}, {"HND", "IAD"}, {"JFK", "HND"}})
	if output[0] != "ICN" || output[1] != "JFK" || output[2] != "HND" || output[3] != "IAD" {
		t.Error("Wrong")
	}
}
