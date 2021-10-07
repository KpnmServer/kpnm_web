
package kweb_util

import (
	bytes   "bytes"
	debug   "runtime/debug"
)


func GetStacks(_n ...int)(stacks []string){
	var n int = 0
	if len(_n) > 0 {
		n = _n[0]
	}

	buf := bytes.NewBuffer(debug.Stack())

	stacks = make([]string, 0)
	for buf.Len() > 0 {
		line, _ := buf.ReadString('\n')
		if len(line) > 1 {
			stacks = append(stacks, line[:len(line)-1])
		}
	}
	if n < -2 || n * 2 + 5 >= len(stacks) {
		return stacks[0:1]
	}
	stacks = append(stacks[0:1], stacks[n * 2 + 5:]...)
	return stacks
}

func GetStack(_n ...int)(string){
	var n int = 0
	if len(_n) > 0 {
		n = _n[0]
	}

	stacks := GetStacks(n + 1)

	buf := bytes.NewBuffer(make([]byte, 0))
	for i, _ := range stacks {
		buf.WriteString(stacks[i])
		buf.WriteByte('\n')
	}
	return buf.String()
}
