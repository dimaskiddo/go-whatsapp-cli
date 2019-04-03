package helper

import (
	"strings"
)

func SplitAtChar(s string, sep string, n int) []interface{} {
	var j int

	var ret []interface{}
	ret = append(ret, "")

	lines := strings.Split(s, sep)
	for i := 0; i < len(lines); i++ {
		if len(ret[j].(string)+sep+lines[i]) > n {
			ret = append(ret, "")
			j++
		}

		ret[j] = ret[j].(string) + sep + lines[i]
	}

	return ret
}
