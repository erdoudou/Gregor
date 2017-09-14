// db
package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var _isInit bool = false
var _dbType string
var _userName string
var _psw string
var _dbName string
var _dbIpAddr string

var _db *sql.DB
var db2 *sql.DB

var _isRelease bool = false

//var _db = &sql.DB{}

func Init(dbType, userName, psw, dbName, dbIpAddr string) {
	_dbType = dbType
	_userName = userName
	_psw = psw
	_dbName = dbName
	_dbIpAddr = dbIpAddr
	_isInit = true

}

func getConnector() *sql.DB {
	if false == _isInit {
		return nil
	}
	if _db == nil {
		if false == connect() {
			return nil
		}
	}
	return _db
}

func connect() bool {
	var connectStr string
	if _isRelease == true {
		connectStr = fmt.Sprintf("%s:%s@/%s?charset=utf8mb4", _userName, _psw, _dbName) //utf8mb4

	} else {
		connectStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", _userName, _psw, _dbIpAddr, _dbName)
	}

	db, err := sql.Open(_dbType, connectStr)
	if nil != err {
		fmt.Println("连接数据库失败")
		return false
	}
	fmt.Println("数据库已连接")
	_db = db
	return true
}

func Query(rowChan chan *sql.Rows, sql string, args ...interface{}) {
	db := getConnector()
	if db == nil {
		rowChan <- nil
		return
	}
	rows, err := db.Query(sql, args...)
	if nil != err {
		fmt.Printf("查询数据库失败:%s\r\n", sql)
		fmt.Println(err.Error())
		rowChan <- nil
	} else {
		rowChan <- rows
	}
}

// 新增查询接口[以返回值得形式,返回结果]
func QueryNormal(sql string, args ...interface{}) (*sql.Rows, bool) {
	db := getConnector()
	if db == nil {
		return nil, false
	}
	rows, err := db.Query(sql, args...)
	if nil != err {
		fmt.Printf("查询数据库失败:%s\r\n", sql)
		fmt.Println(err.Error())
		return nil, false
	} else {
		return rows, true
	}
}

func ExecInsert(idChan chan int, rstChan chan bool, sql string, args ...interface{}) {
	db := getConnector()
	if db == nil {
		fmt.Println("无法连接到数据库")
		idChan <- 0
		rstChan <- false
		return
	}
	rst, err := db.Exec(sql, args...)
	if nil != err {
		fmt.Printf("插入数据库失败:%s\r\n", sql)
		fmt.Println(err.Error())
		idChan <- 0
		rstChan <- false
		return
	}
	updateNum, err := rst.RowsAffected()
	if nil != err || updateNum <= 0 {
		fmt.Println("没有数据被插入")
		idChan <- 0
		rstChan <- false
		return
	}
	insertID, err := rst.LastInsertId()
	if nil != err {
		fmt.Println("未能获取到自增长ID")
		idChan <- 0
		rstChan <- true
		return
	}
	idChan <- int(insertID)
	rstChan <- true
	return
}

func ExecUpdate(updateNumChan chan int, rstChan chan bool, sql string, args ...interface{}) {
	db := getConnector()
	if db == nil {
		updateNumChan <- 0
		rstChan <- false
		return
	}
	rst, err := db.Exec(sql, args...)
	if nil != err {
		fmt.Println("插入数据库失败:", sql, ";", args)
		fmt.Println(err.Error())
		updateNumChan <- 0
		rstChan <- false
		return
	}
	updateNum, err := rst.RowsAffected()
	updateNumChan <- int(updateNum)
	rstChan <- err == nil
}

// add by hjd 20170301
// 执行插入操作,不返回自增长ID
func ExecInsertWithoutInsertId(idChan chan int, rstChan chan bool, sql string, args ...interface{}) {
	db := getConnector()
	if db == nil {
		fmt.Println("无法连接到数据库")
		idChan <- 0
		rstChan <- false
		return
	}
	rst, err := db.Exec(sql, args...)
	if nil != err {
		fmt.Printf("插入数据库失败:%s\r\n", sql)
		fmt.Println(err.Error())
		idChan <- 0
		rstChan <- false
		return
	}
	updateNum, err := rst.RowsAffected()
	if nil != err || updateNum <= 0 {
		fmt.Println("没有数据被插入")
		idChan <- 0
		rstChan <- false
		return
	}

	rstChan <- true
	return
}

func ExecInsertAccount(idChan chan int, rstChan chan bool, sql string, userName string, password string) {
	//db := getConnector()
	//if db == nil {
	//fmt.Println("无法连接到数据库")
	//idChan <- 0
	//rstChan <- false
	//return
	//}
	//rst, err := db2.Exec("insert into account (username,password) values (?,?)", "hejin", "user")
	start := time.Now()
	//for i := 1001; i <= 1100; i++ {
	//每次循环内部都会去连接池获取一个新的连接，效率低下
	//db2.Exec("insert into account (username,password) values (?,?)", "hejin", "user") //+strconv.Itoa(i)
	//}
	//end := time.Now()
	//fmt.Println("方式1 insert total time:", end.Sub(start).Seconds())
	//return
	_, _ = db2.Exec(sql, userName, password)
	end := time.Now()
	fmt.Println("方式1 insert total time:", end.Sub(start).Seconds())
}

func Init2() {
	db2, _ = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", "root", "root", "192.168.3.52:3306", "gun"))
}

func Insert() {

	//方式1 insert
	//strconv,int转string:strconv.Itoa(i)
	start := time.Now()
	for i := 1001; i <= 1100; i++ {
		//每次循环内部都会去连接池获取一个新的连接，效率低下
		db2.Exec("insert into account (username,password) values (?,?)", "hejin", "user"+strconv.Itoa(i))
	}
	end := time.Now()
	fmt.Println("方式1 insert total time:", end.Sub(start).Seconds())

	//方式2 insert
	start = time.Now()
	for i := 1101; i <= 1200; i++ {
		//Prepare函数每次循环内部都会去连接池获取一个新的连接，效率低下
		stm, _ := db2.Prepare("insert into account (username,password) values (?,?)")
		stm.Exec("hejin", "user"+strconv.Itoa(i))
		stm.Close()
	}
	end = time.Now()
	fmt.Println("方式2 insert total time:", end.Sub(start).Seconds())

	//方式3 insert
	start = time.Now()
	stm, _ := db2.Prepare("insert into account (username,password) values (?,?)")
	for i := 1201; i <= 1300; i++ {
		//Exec内部并没有去获取连接，为什么效率还是低呢？
		stm.Exec("hejin", "user"+strconv.Itoa(i))
	}
	stm.Close()
	end = time.Now()
	fmt.Println("方式3 insert total time:", end.Sub(start).Seconds())

	//方式4 insert
	start = time.Now()
	//Begin函数内部会去获取连接
	tx, _ := db2.Begin()
	for i := 1301; i <= 1400; i++ {
		//每次循环用的都是tx内部的连接，没有新建连接，效率高
		tx.Exec("insert into account (username,password) values (?,?)", "hejin", "user"+strconv.Itoa(i))
	}
	//最后释放tx内部的连接
	tx.Commit()

	end = time.Now()
	fmt.Println("方式4 insert total time:", end.Sub(start).Seconds())

}

// add end
