package libanreport

import "bytes"

type DataJSON struct {
	Type int
	Key  string
	Val  string
}

var noBreakSpace = '\u00A0'

func removeSpecialFromDataJSON(refDatas []DataJSON) {
	max := len(refDatas)
	for i := 0; i < max; i++ {
		refDatas[i].Val = removeSpecialRune(refDatas[i].Val)
	}
}

func removeSpecialRune(txt string) string {
	var buff bytes.Buffer
	for _, r := range txt {
		if r == noBreakSpace {
			continue
		}
		buff.WriteRune(r)
	}
	return buff.String()
}
