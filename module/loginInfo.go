package module

func GetDashBoard(adminname, token string, adminLevel int) (data map[string]interface{}) {
	data = map[string]interface{}{}
	data["adminname"] = adminname
	data["token"] = token
	data["admin_level"] = adminLevel
	permissionList := []interface{}{}
	childrenList := []interface{}{}
	childrenMeta1 := MakeMeta("dashboard", "DROP_OS")
	childrenMap := MakeMap2("dashboard", "Dashboard", childrenMeta1)
	childrenList = append(childrenList, childrenMap)
	map1 := MakeMap("/", "Layout", "/dashboard", childrenList)
	permissionList = append(permissionList, map1)
	childrenMeta2 := MakeMeta("el-icon-s-help", "后台管理")
	childrenList2 := []interface{}{}
	childrenMeta3 := MakeMeta("x", "成员考勤")
	childrenMap2 := MakeMap2("table", "Table", childrenMeta3)
	childrenList2 = append(childrenList2, childrenMap2)
	childrenMeta4 := MakeMeta("x", "成员列表")
	childrenMap3 := MakeMap2("table2", "Table2", childrenMeta4)
	childrenList2 = append(childrenList2, childrenMap3)
	//childrenMeta6 := MakeMeta("x", "发送邮箱")
	//childrenMap5 := MakeMap2("send", "Send", childrenMeta6)
	//childrenList2 = append(childrenList2, childrenMap5)
	childrenMeta7 := MakeMeta("x", "成员签到详情")
	childrenMap6 := MakeMapHidden("details", "Details", childrenMeta7)
	childrenList2 = append(childrenList2, childrenMap6)
	childrenMeta12 := MakeMeta("x", "范围查询")
	childrenMap11 := MakeMapHidden("echarts", "Echarts", childrenMeta12)
	childrenList2 = append(childrenList2, childrenMap11)
	childrenMeta8 := MakeMeta("x", "成员签到详情echarts版本")
	childrenMap7 := MakeMapHidden("detailsEcharts", "detailsEcharts", childrenMeta8)
	childrenList2 = append(childrenList2, childrenMap7)
	childrenMeta9 := MakeMeta("x", "成员签到详情")
	childrenMap8 := MakeMapHidden("perInfo", "perInfo", childrenMeta9)
	childrenList2 = append(childrenList2, childrenMap8)
	childrenMeta10 := MakeMeta("x", "可疑数据")
	childrenMap9 := MakeMap2("warnInfo", "warnInfo", childrenMeta10)
	childrenList2 = append(childrenList2, childrenMap9)
	if adminLevel == 0{
		childrenMeta5 := MakeMeta("x", "申请列表")
		childrenMap4 := MakeMap2("apply", "Apply", childrenMeta5)
		childrenList2 = append(childrenList2, childrenMap4)
		childrenMeta11 := MakeMeta("x", "添加管理员")
		childrenMap10 := MakeMap2("addAdmin", "addAdmin", childrenMeta11)
		childrenList2 = append(childrenList2, childrenMap10)
	}
	map2 := MakeMap1("/student", "Layout", "/student/table", "Student", childrenMeta2, childrenList2)
	permissionList = append(permissionList, map2)
	data["permission"] = permissionList
	return data
}

func MakeMap(path, component, redirect string, children []interface{}) (data map[string]interface{}) {
	data = map[string]interface{}{}
	data["path"] = path
	data["component"] = component
	data["redirect"] = redirect
	data["children"] = children
	return data
}

func MakeMap1(path, component, redirect, name string, meta map[string]string, children []interface{}) (data map[string]interface{}) {
	data = map[string]interface{}{}
	data["path"] = path
	data["component"] = component
	data["redirect"] = redirect
	data["name"] = name
	data["meta"] = meta
	data["children"] = children
	return data
}

func MakeMap2(path, name string, meta map[string]string) (data map[string]interface{}) {
	data = map[string]interface{}{}
	data["path"] = path
	data["name"] = name
	data["meta"] = meta
	return data
}

func MakeMapHidden(path, name string, meta map[string]string) (data map[string]interface{}) {
	data = map[string]interface{}{}
	data["path"] = path
	data["name"] = name
	data["meta"] = meta
	data["hidden"] = "true"
	return data
}

func MakeMeta(icon, title string) (meta map[string]string) {
	meta = map[string]string{}
	meta["icon"] = icon
	meta["title"] = title
	return meta
}
