package follow_set

import (
	"strings"
	"fmt"
	"compiler-LL1/rule"
	"compiler-LL1/first_set"
	"compiler-LL1/util"
)

type FollowSet map[byte]map[byte]struct{}

func (f FollowSet) String() string {
	var build strings.Builder
	for key, value := range f {
		build.WriteString(fmt.Sprintf("FOLLOW(%c)={ ", key))
		for item := range value {
			build.WriteString(fmt.Sprintf("%c ",item))
		}
		build.WriteString("}\n")
	}
	return build.String()
}

func (f FollowSet) Strings() []string {
	var build strings.Builder
	var ans []string
	for key, value := range f {
		build.WriteString(fmt.Sprintf("FOLLOW(%c)={ ", key))
		for item := range value {
			build.WriteString(fmt.Sprintf("%c ", item))
		}
		build.WriteString("}")
		ans = append(ans, build.String())
		build.Reset()
	}
	return ans
}

func GetFollowFrom(rules rule.Rules,start byte,firstSet first_set.FirstSet) FollowSet {
	followSet := make(FollowSet)

	followSet[start] = make(map[byte]struct{})
	followSet[start]['#'] = struct{}{}

	var changed bool
	for {
		changed = false
		for key, r := range rules {
			for _, v := range r {
				for i := 0; i < len(v)-1; i++ {
					// A->aBb
					if followSet[v[i]] == nil {
						followSet[v[i]] = make(map[byte]struct{})
					}
					if !util.IsTerminal(v[i]) {
						if util.IsTerminal(v[i+1]) {
							if MergeSet(followSet[v[i]], map[byte]struct{}{v[i+1]: {}}) != 0{
								changed = true
							}
						} else {
							if RemoveEmptyAndMergeSet(followSet[v[i]], firstSet[v[i+1]]) != 0{
								changed = true
							}
						}
						if firstSet.HaveEmpty(v[i+1]) {
							if MergeSet(followSet[v[i]], followSet[key]) != 0{
								changed = true
							}
						}
					}
				}

				// A->aB
				if followSet[v[len(v)-1]] == nil {
					followSet[v[len(v)-1]] = make(map[byte]struct{})
				}
				if !util.IsTerminal(v[len(v)-1]){
					if MergeSet(followSet[v[len(v)-1]], followSet[key]) != 0{
						changed = true
					}
				}
			}
		}
		if !changed {
			break
		}
	}
	return followSet
}

func RemoveEmptyAndMergeSet(a map[byte]struct{}, b map[byte]struct{}) int {
	flag := false
	if _, flag = b['@'];flag {
		flag = true
		delete(b,'@')
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
