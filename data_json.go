package libanreport

const WrapTextTypeNone = 0b00
const WrapTextTypeNewLine = 0b01

type DataJSON struct {
	Type         int
	Key          string
	Val          string
	WrapTextType int
}

//func removeSpecialFromDataJSON(refDatas []DataJSON) {
//	max := len(refDatas)
//	for i := 0; i < max; i++ {
//		refDatas[i].Val = removeSpecialRune(refDatas[i].Val)
//	}
//}
