package helper

import (
	"strings"
)

func SplitAfterCharN(s string, sep string, lim int, n int, pre bool) []string {
	ret := []string{""}
	s = strings.TrimSpace(s)

	if len(s) > 0 {
		var j int

		add := ""
		if pre {
			add = "```"
			lim = lim - 6
		}
		ret[0] = add

		split := strings.SplitN(s, sep, n)
		for i := 0; i < len(split); i++ {
			if len(ret[j]+sep+split[i]) > lim {
				ret[j] = strings.TrimSpace(ret[j] + add)

				ret = append(ret, add)
				j++
			}

			if len(ret[j])-len(add) == 0 {
				ret[j] = ret[j] + split[i]
			} else {
				ret[j] = ret[j] + sep + split[i]
			}

			if i == len(split)-1 {
				ret[j] = strings.TrimSpace(ret[j] + add)
			}
		}
	}

	return ret
}

func SplitWithEscapeN(s string, sep string, n int) []string {
	ret := []string{""}
	s = strings.TrimSpace(s)

	if len(s) > 0 {
		var j int

		split := strings.SplitN(s, sep, n)
		for i := 0; i < len(split); i++ {
			ret[j] = strings.TrimSpace(split[i])
			for (strings.Count(ret[j], "'") == 1 || strings.Count(ret[j], "\"") == 1) && i != len(split)-1 {
				ret[j] = ret[j] + sep + strings.TrimSpace(split[i+1])
				i++
			}

			if i != len(split)-1 {
				ret = append(ret, "")
				j++
			}
		}
	}

	return ret
}
