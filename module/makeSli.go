package module

// 返回列表的切片
func MakeSli(pageNum interface{}, list []interface{}) []interface{} {
	pagenum2 := pageNum.(float64)
	pagenum3 := int(pagenum2)
	pageNumMax := len(list) / 10
	if pagenum3 <= pageNumMax {
		return list[pagenum3*10-10 : pagenum3*10]
	}
	if pagenum3 == pageNumMax+1 {
		rightNum := len(list) - len(list)/10*10
		return list[pagenum3*10-10 : pagenum3*10-10+rightNum]
	}
	return nil
}
