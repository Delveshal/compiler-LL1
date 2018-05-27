package rule

import (
	"fmt"
	"strings"
	"github.com/Delveshal/compiler-LL1/util"
)

type Rules map[byte][]string

func NewRules() Rules {
	return make(Rules)
}

func (r Rules) AddRules(s string) error {
	stepOne := strings.Split(s, "->")
	if len(stepOne) != 2 {
		return fmt.Errorf(fmt.Sprintf("the format of input is invalid,expect X->Y but actually %s", s))
	}
	if len(stepOne[0]) != 1 {
		return fmt.Errorf(fmt.Sprintf("input is invalid,expect X on the left but actually %s", stepOne[0]))
	}
	stepTwo := strings.Split(strings.Replace(stepOne[1], " ", "", -1), "|")
	for i := range stepTwo {
		r[stepOne[0][0]] = append(r[stepOne[0][0]], stepTwo[i])
	}
	return nil
}

func (r Rules) String() string {
	var builder strings.Builder
	for fist, second := range r {
		for i := range second {
			builder.WriteString(fmt.Sprintf("%c->%s\n", fist, second[i]))
		}
	}
	return builder.String()
}

func AllIsTer(s string) bool {
	for i := range s {
		if s[i] < 'A' || s[i] > 'Z' {
			return false
		}
	}
	return true
}

// 判断是否存在有 X->e
func (r Rules) HaveEmptyFormula(first byte) bool {
	for _, value := range r[first] {
		if value == "@" {
			return true
		}
	}
	return false
}

func (r Rules) TheFirstItemIs(first, item byte) string {
	for _, value := range r[first] {
		if value[0] == item {
			return value
		}
	}
	return ""
}

func (r Rules) Dfs(first, terminal byte) bool {
	for i := range r[first] {
		if util.IsTerminal(r[first][i][0]) && r[first][i][0] == terminal {
			return true
		} else if !util.IsTerminal(r[first][i][0]) {
			if ok := r.Dfs(r[first][i][0], terminal); ok {
				return true
			}
		}
	}
	return false
}
