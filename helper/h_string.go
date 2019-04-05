package helper

import (
	"strings"
)

func SplitAtChar(s string, sep string, n int, out string) []string {
	var j int

	var strPre string
	strPre = ""

	if out == "pretty" {
		strPre = "```"
		n = n - 6
	}

	var strRow []string
	strRow = append(strRow, "")

	rows := strings.Split(s, sep)
	for i := 0; i < len(rows); i++ {
		if len(strRow[j]+sep+rows[i]) > n || i == len(rows)-1 {
			strRow[j] = strings.TrimSpace(strRow[j]) + strPre

			if i != len(rows)-1 {
				strRow = append(strRow, "")
				j++
			}
		}

		if len(strRow[j]) == 0 {
			strRow[j] = strPre + rows[i]
		} else {
			strRow[j] = strRow[j] + sep + rows[i]
		}
	}

	return strRow
}
