package libanreport

type DataJSON struct {
	Type int
	Key  string
	Val  string
}

//func removeSpecialFromDataJSON(refDatas []DataJSON) {
//	max := len(refDatas)
//	for i := 0; i < max; i++ {
//		refDatas[i].Val = removeSpecialRune(refDatas[i].Val)
//	}
//}
