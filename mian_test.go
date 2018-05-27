package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"github.com/Delveshal/compiler-LL1/rule"
	"github.com/Delveshal/compiler-LL1/first_set"
	"github.com/Delveshal/compiler-LL1/follow_set"
	"github.com/Delveshal/compiler-LL1/chart"
	"github.com/Delveshal/compiler-LL1/analyzer"
)

type testCase struct {
	id        int
	Rule      rule.Rules           `json:"rule"`
	Input     string               `json:"input"`
	FirstSet  first_set.FirstSet   `json:"first_set"`
	FollowSet follow_set.FollowSet `json:"follow_set"`
	Chart     chart.Chart          `json:"chart"`
	Step      []*analyzer.Item     `json:"step"`
}

func TestAll(t *testing.T) {
	dir := "test_data"
	i := 0
	for {
		if !PathExist(dir + string(filepath.Separator) + strconv.Itoa(i) + ".in") {
			break
		}
		in, err := os.Open(dir + string(filepath.Separator) + strconv.Itoa(i) + ".in")
		if err != nil {
			t.Error(err)
		}
		out, err := os.Open(dir + string(filepath.Separator) + strconv.Itoa(i) + ".out")
		if err != nil {
			t.Error(err)
		}
		bufIn, err := ioutil.ReadAll(in)
		if err != nil {
			t.Error(err)
		}
		bufOut, err := ioutil.ReadAll(out)
		if err != nil {
			t.Error(err)
		}
		err = Check(i, bufIn, bufOut)
		if err != nil {
			t.Error(err.Error())
		}
		i++
	}
}

func Check(id int, bufIn, bufOut []byte) error {
	t := bytes.IndexByte(bufIn, '#')
	raw := string(bufIn[:t])

	grammar := strings.Split(raw, "\n")
	rules := rule.NewRules()
	for i := range grammar {
		rules.AddRules(grammar[i])
	}
	firstSet := first_set.GetFirstFrom(rules)
	start := grammar[0][0]
	followSet := follow_set.GetFollowFrom(rules, start, firstSet)
	ch := chart.GetChartFrom(firstSet, followSet, rules)
	input := strings.Replace(string(bufIn[t+1:]), "\n", "", -1)
	step, err := analyzer.Analyze(ch, start, input)
	if err != nil {
		return err
	}
	actually := &testCase{
		id:        id,
		Rule:      rules,
		FirstSet:  firstSet,
		FollowSet: followSet,
		Chart:     ch,
		Step:      step,
	}
	expect := &testCase{
		id: id,
	}
	err = json.Unmarshal(bufOut, expect)
	if err != nil {
		return err
	}
	if ok := reflect.DeepEqual(actually, expect); !ok {
		return fmt.Errorf("not equal")
	}
	return nil
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
