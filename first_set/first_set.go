package first_set

import (
	"fmt"
	"strings"
	"github.com/Delveshal/compiler-LL1/rule"
)

type FirstSet map[byte]map[byte]struct{}

func GetFirstFrom(rules rule.Rules) FirstSet {
	firstSet := make(FirstSet)
	var changed bool
	for {
		changed = false
		for key, r := range rules {
			if firstSet[key] == nil {
				firstSet[key] = make(map[byte]struct{})
			}
			for _, v := range r {
				// 第一个是非终结符
				if v[0] < 'A' || v[0] > 'Z' {
					if MergeSet(firstSet[key], map[byte]struct{}{v[0]: {}}) != 0 {
						changed = true
					}
					continue
				}
				// 第一个是非终结符
				if v[0] >= 'A' && v[0] <= 'Z' {
					if RemoveEmptyAndMergeSet(firstSet[key], firstSet[v[0]]) != 0 {
						changed = true
					}
				}
				// 查看后面是不是终结符
				if rule.AllIsTer(v) {
					var j = 0
					for j < len(v) {
						// 含有empty表达式把其first合到当前first集
						if rules.HaveEmptyFormula(v[j]) {
							if RemoveEmptyAndMergeSet(firstSet[key], firstSet[v[j]]) != 0 {
								changed = true
							}
							j++
						} else {
							break
						}
					}
					// 全部都含有empty表达式
					if j == len(v) {
						if MergeSet(firstSet[key], map[byte]struct{}{'@': {}}) != 0 {
							changed = true
						}
					}
				}
			}
		}
		if !changed {
			break
		}
	}
	return firstSet
}

func RemoveEmptyAndMergeSet(a map[byte]struct{}, b map[byte]struct{}) int {
	flag := false
	if _, flag = b['@']; flag {
		flag = true
		delete(b, '@')
	}
	count := 0
	for key, value := range b {
		if _, ok := a[key]; !ok {
			count++
		}
		a[key] = value
	}
	if flag {
		b['@'] = struct{}{}
	}
	return count
}

func MergeSet(a map[byte]struct{}, b map[byte]struct{}) int {
	count := 0
	for key, value := range b {
		if _, ok := a[key]; !ok {
			count++
		}
		a[key] = value
	}
	return count
}

func (f FirstSet) String() string {
	var build strings.Builder
	for key, value := range f {
		build.WriteString(fmt.Sprintf("FIRST(%c)={ ", key))
		for item := range value {
			build.WriteString(fmt.Sprintf("%c ", item))
		}
		build.WriteString("}\n")
	}
	return build.String()
}

func (f FirstSet) Strings() []string {
	var build strings.Builder
	var ans []string
	for key, value := range f {
		build.WriteString(fmt.Sprintf("FIRST(%c)={ ", key))
		for item := range value {
			build.WriteString(fmt.Sprintf("%c ", item))
		}
		build.WriteString("}")
		ans = append(ans, build.String())
		build.Reset()
	}
	return ans
}

func (f FirstSet) HaveEmpty(first byte) bool {
	_, ok := f[first]['@']
	return ok
}
