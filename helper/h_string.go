package helper

import (
	"strings"
)

func SplitAtChar(s string, sep string, n int) []string {
	var j int
	var row string

	var strRow []string
	strRow = append(strRow, "")

	rows := strings.Split(s, sep)
	for i := 0; i < len(rows); i++ {
		switch len(strRow[j]) {
		case 0:
			row = strRow[j] + rows[i]
		default:
			row = strRow[j] + sep + rows[i]
		}

		if len(row) > n {
			strRow[j] = strings.TrimSpace(strRow[j])
			strRow = append(strRow, "")
			j++
		}

		strRow[j] = row
	}

	return strRow
}
