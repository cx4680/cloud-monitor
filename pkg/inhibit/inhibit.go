package inhibit

import (
	"sort"
	"strconv"
	"strings"
)

type InhibitRule struct {
	SourceMatchers []Matcher `json:"sourceMatchers"`
	TargetMatchers []Matcher `json:"targetMatchers"`
	Equal          []string  `json:"equal"`
}

type Matcher struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	MatchType string `json:"matchType"`
	Regex     bool   `json:"regex"`
}

func BuildRules(levels []uint8) []InhibitRule {
	var rules []InhibitRule
	sort.Slice(levels, func(i, j int) bool {
		return levels[i] < levels[j]
	})

	for i, level := range levels {
		if i == len(levels)-1 {
			break
		}
		rule := InhibitRule{
			SourceMatchers: []Matcher{{
				Name:  "app",
				Value: "hawkeye",
				Regex: false,
			}, {
				Name:  "severity",
				Value: strconv.Itoa(int(level)),
				Regex: false,
			}},
			TargetMatchers: []Matcher{{
				Name:  "app",
				Value: "hawkeye",
				Regex: false,
			}, {
				Name:  "severity",
				Value: strings.Join(getGreaterLevels(level, levels), "|"),
				Regex: true,
			}},
			Equal: []string{
				"alertname",
			},
		}
		rules = append(rules, rule)
	}

	return rules
}

func getGreaterLevels(level uint8, levels []uint8) []string {
	var greaters []string
	for _, l := range levels {
		if l > level {
			greaters = append(greaters, strconv.Itoa(int(l)))
		}
	}
	return greaters
}
