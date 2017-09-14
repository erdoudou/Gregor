package roleMgr

import (
	"Public/db"

	//"database/sql"

	"fmt"
)

var _roleMap map[string]*Role = make(map[string]*Role)

type Role struct {
	RoleName string
	RolePwd  string
}

//玩家的名字可以作为唯一的标识

//从数据库中获得所有的玩家
func GetAllRole() {
	row, ret := db.QueryNormal("select username,password from account")
	if ret == false {
		fmt.Println("查询错误")
		return
	}

	if row == nil {
		fmt.Println("数据错误")
		return
	}

	defer row.Close()

	var roleName, pwd string

	for row.Next() != false {

		err := row.Scan(&roleName, &pwd)

		if err != nil {
			fmt.Println("解析失败")
			return
		}
		role := &Role{
			RoleName: roleName,
			RolePwd:  pwd,
		}
		_roleMap[roleName] = role

	}
}

//判断一个玩家是否存在数据库中
func IsRole(name string) bool {
	if _, ok := _roleMap[name]; ok == false {
		fmt.Println("玩家不存在")
		return false
	}

	return true

}

//判断一个玩家密码是否正确
func IsPwd(name, pwd string) bool {
	if _roleMap[name].RolePwd == pwd {
		return true
	}

	return false

}

//将玩家存入数据库中
func SetRole(name, pwd string) {
	role := &Role{
		RoleName: name,
		RolePwd:  pwd,
	}

	_roleMap[name] = role

	idchan := make(chan int, 2)
	reschan := make(chan bool, 2)

	insertCmd := "INSERT  INTO account(username,password) VALUE (?,?)"
	db.ExecInsert(idchan, reschan, insertCmd, name, pwd)

}
