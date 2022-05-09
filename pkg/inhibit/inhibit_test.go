package inhibit

import (
	"fmt"
	"sort"
	"testing"
)

func Test_sortLevels(t *testing.T) {
	levels := []uint8{1, 2, 5, 4, 3, 0}
	sort.Slice(levels, func(i, j int) bool {
		return levels[i] < levels[j]
	})

	fmt.Println(levels)
}

func Test_BuildRules(t *testing.T) {
	levels := []uint8{1, 2, 4, 3}
	rules := BuildRules(levels)
	fmt.Println(rules)
}
