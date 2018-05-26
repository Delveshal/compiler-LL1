package chart

import (
	"strings"
	"compiler-LL1/first_set"
	"compiler-LL1/follow_set"
	"compiler-LL1/rule"
)

type Chart map[byte]map[byte]string

func GetChartFrom(first first_set.FirstSet, follow follow_set.FollowSet, rules rule.Rules) Chart {
	chart := make(Chart)
	for A, r := range rules{
		for i := range r {
			if chart[A] == nil {
				chart[A] = make(map[byte]string)
			}
			for a := range first[A] {
				if a != '@' {
					if t := rules.TheFirstItemIs(A, a); t != "" {
						chart[A][a] = string(A) + "->" + t
					} else if rules.Dfs(A, a) {
						chart[A][a] = string(A) + "->" + r[i]
					}
				}else {
					for k := range follow[A]{
						chart[A][k] = string(A) + "->" + string(a)
					}
				}
			}
		}
	}
	return chart
}

func (c Chart) String() string {
	var builder strings.Builder
	for row, v := range c {
		for col, formula := range v {
			builder.Write([]byte{'{', '[', row, ' ', col, ']', ' '})
			builder.WriteString(formula + "} ")
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}

func (c Chart) CoverToTable() [][]string {
	terminal := 1
	notTerminal := 1
	index := make(map[byte]int)
	for key,r := range c{
		if index[key] == 0{
			index[key] = notTerminal
			notTerminal++
		}
		for k := range r{
			if index[k] == 0{
				index[k] = terminal
				terminal++
			}
		}
	}
	result := make([][]string,notTerminal)
	for i := range result{
		result[i] = make([]string,terminal)
	}
	for key,r := range c{
		result[index[key]][0] = string(key)
		for k := range r {
			result[index[key]][index[k]] = c[key][k]
			result[0][index[k]] = string(k)
		}
	}
	return result
}