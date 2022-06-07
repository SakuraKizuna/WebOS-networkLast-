package models

import "drop_os_back/dao"

type Todo struct {
	Id       int    `json:"id" gorm:"primary_key"`
	Content  string `json:"content"`
	Username string `json:"username"`
}

func GetTodoList() (todos []*Todo, err error) {
	var todoList []*Todo
	//err = dao.DB.Find(&todoList).Error
	sql := "SELECT * FROM `todos` ORDER BY id desc"
	err = dao.DB.Raw(sql).Scan(&todoList).Error
	if err != nil {
		return nil, err
	}
	return todoList, err
}

func DealTodoList(todoList []*Todo) (todosli []interface{}) {
	TodoSLi := []interface{}{}
	for _, v := range todoList {
		todoMap := map[string]interface{}{}
		todoMap["id"] = v.Id
		todoMap["content"] = v.Content
		TodoSLi = append(TodoSLi, todoMap)
	}
	return TodoSLi
}

func DeleteTodo(id int) error {
	todoInfo := Todo{}
	err := dao.DB.Where("id=?", id).Delete(&todoInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func AddTodo(username, content string) error {
	todoInfo := Todo{}
	todoInfo.Content = content
	todoInfo.Username = username
	err := dao.DB.Save(&todoInfo).Error
	if err != nil {
		return err
	}
	return nil
}
