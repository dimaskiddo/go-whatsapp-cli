package hlp

import (
	"strings"
)

func SplitAfterCharN(s string, sep string, lim int, n int, pre bool, trim bool) []string {
	ret := []string{""}

	if trim {
		s = strings.TrimSpace(s)
	}

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
				if trim {
					ret[j] = strings.TrimSpace(ret[j] + add)
				} else {
					ret[j] = ret[j] + add
				}

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

func SplitWithEscapeN(s string, sep string, n int, trim bool) []string {
	ret := []string{""}

	if trim {
		s = strings.TrimSpace(s)
	}

	if len(s) > 0 {
		var j int

		split := strings.SplitN(s, sep, n)
		for i := 0; i < len(split); i++ {
			if trim {
				ret[j] = strings.TrimSpace(split[i])
			} else {
				ret[j] = split[i]
			}

			if (strings.Count(ret[j], "'") == 1 || strings.Count(ret[j], "\"") == 1) && i != len(split)-1 {
				if trim {
					ret[j] = ret[j] + sep + strings.TrimSpace(split[i+1])
				} else {
					ret[j] = ret[j] + sep + split[i+1]
				}
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
