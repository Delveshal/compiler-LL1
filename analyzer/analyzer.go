package analyzer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Delveshal/compiler-LL1/chart"
	"github.com/Delveshal/compiler-LL1/util"
)

type Item struct {
	Symbol string `json:"symbol"`
	Cur    string `json:"cur"`
	Input  string `json:"input"`
	Mark   string `json:"mark"`
}

func Analyze(chart chart.Chart, start byte, input string) ([]*Item, error) {
	if len(input) <= 0 {
		return nil, errors.New("empty")
	}
	if input[len(input)-1] != '#' {
		input += "#"
	}
	symbol := []byte{'#', start}
	var step []*Item
	a := input[0]
	input = input[1:]
	for {
		step = append(step, &Item{
			Symbol: string(symbol),
			Cur:    string(a),
			Input:  input,
		})
		x := symbol[len(symbol)-1]
		symbol = symbol[:len(symbol)-1]
		if util.IsTerminal(x) {
			if x == a {
				if a != '#' {
					a = input[0]
					input = input[1:]
				}
			} else {
				return nil, errors.New(fmt.Sprintf("%c != %c", x, a))
			}
		} else {
			if t, ok := chart[x][a]; ok {
				index := bytes.LastIndexByte([]byte(t), '>')
				if t[len(t)-1] != '@' {
					for i := len(t) - 1; i > index; i-- {
						symbol = append(symbol, t[i])
					}
				}
			} else {
				return nil, errors.New("the shell is empty: " + fmt.Sprintf("%c %c", x, a))
			}
		}
		if x == '#' {
			break
		}
	}
	return step, nil
}
