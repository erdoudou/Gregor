package box

import (
	"Public/db"
	"com/ViKing/FPSDemo/hallServer/overallsituation"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

/*背包设计背包的mgr*/
// 拥有随机宝箱最大数
//var MAX_BOX_SLOTS int = overallsituation.OverallData.MaxBoxSlots

type TreasureHouse struct {
	boxContent string
	boxNum     int
}

//获得一个玩家身上的所有宝箱
func (treasureHouse *TreasureHouse) LoadAllRoleBox(roleid int64) {

	row1 := make(chan *sql.Rows, 1)
	db.Query(row1, "SELECT treasureboxs FROM treasurehouse WHERE useruid = ?",
		roleid)
	in1 := <-row1

	if in1 == nil {
		return
	}

	defer in1.Close()

	var treasureboxs string = ""

	if in1.Next() != false {

		err := in1.Scan(&treasureboxs)

		if err != nil {
			return
		}

	}

	treasureHouse.Init(treasureboxs)

	return

}

//初始化玩家身上的宝箱数据
func (treasureHouse *TreasureHouse) Init(content string) {
	treasureHouse.boxContent = content
	treasureHouse.CalculateBoxNum()
}

// 获取宝箱数量
func (treasureHouse *TreasureHouse) GetBoxNum() int {
	return treasureHouse.boxNum
}

// 获取宝箱内容
func (treasureHouse *TreasureHouse) GetBoxContent() string {
	return treasureHouse.boxContent
}

// 计算宝箱数量
func (treasureHouse *TreasureHouse) CalculateBoxNum() {
	if treasureHouse.boxContent == "" {
		treasureHouse.boxNum = 0
		return
	}
	treasureHouse.boxNum = strings.Count(treasureHouse.boxContent, ",")
}

// 玩家身上宝箱数量是否为空
func (treasureHouse *TreasureHouse) IsEmpty() bool {
	if treasureHouse.boxNum <= 0 {
		return true
	} else {
		return false
	}
}

// 判断潘家是否还能获得宝箱
func (treasureHouse *TreasureHouse) HasSlots() bool {
	if overallsituation.OverallData.MaxBoxSlots-treasureHouse.boxNum <= 0 {
		return false
	} else {
		return true
	}
}

func (treasureHouse *TreasureHouse) HasBoxs(boxuid int64) bool {
	uidStr := strconv.Itoa(int(boxuid))
	ret := strings.Contains(treasureHouse.boxContent, uidStr)
	if ret {
		return true
	}
	return false
}

//判断一个玩家身上曾经是否有宝箱数据
func (treasureHouse *TreasureHouse) AddBox(boxUid, roleId int64) bool {

	if treasureHouse.HasSlots() == false {
		return false
	}

	row, ret := db.QueryNormal("SELECT * FROM treasurehouse WHERE useruid = ?", roleId)
	if ret == false {
		fmt.Println("从数据库中取出玩家宝箱数据错误")
		return false
	}

	defer row.Close()

	if row == nil || row.Next() == false {
		//在数据库中新添宝箱数据
		itemStr := strconv.Itoa(int(boxUid)) + ","
		treasureHouse.boxContent += itemStr
		fmt.Println("玩家身上没有新的宝箱数据", treasureHouse.boxContent)
		treasureHouse.boxNum = treasureHouse.boxNum + 1
		treasureHouse.insertHouseData(roleId)
	} else {
		itemStr := strconv.Itoa(int(boxUid)) + ","

		treasureHouse.boxContent += itemStr
		fmt.Println("玩家身上有宝箱数据", treasureHouse.boxContent)
		treasureHouse.boxNum = treasureHouse.boxNum + 1
		treasureHouse.updateHouseData(roleId)
	}

	return true

}

//数据库中背包数据的更新
func (treasureHouse *TreasureHouse) updateHouseData(roleId int64) {

	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE treasurehouse SET treasureboxs=? WHERE useruid=?"

	db.ExecUpdate(idchan1, reschan1, updateCmd,
		treasureHouse.boxContent, roleId)

}

//插入数据
func (treasureHouse *TreasureHouse) insertHouseData(roleId int64) {

	idchan := make(chan int, 2)
	reschan := make(chan bool, 2)

	insertCmd := "INSERT  INTO treasurehouse(useruid,treasureboxs) VALUE (?,?)"
	db.ExecInsert(idchan, reschan, insertCmd, roleId, treasureHouse.boxContent)

}

func (treasureHouse *TreasureHouse) SaveHouseData(roleId int64) {
	row, ret := db.QueryNormal("SELECT * FROM treasurehouse WHERE useruid = ?", roleId)
	if ret == false {
		fmt.Println("从数据库中查找house数据失败")
		return
	}

	defer row.Close()

	if row == nil || row.Next() == false {
		fmt.Println("house中没有数据")
		//在数据库中新添item数据
		treasureHouse.insertHouseData(roleId)

	} else {
		treasureHouse.updateHouseData(roleId)
	}

}

//玩家下线存数据
func insertHouse(roleId int64, boxContent string) {

	idchan := make(chan int, 2)
	reschan := make(chan bool, 2)

	insertCmd := "INSERT  INTO treasurehouse(useruid,treasureboxs) VALUE (?,?)"
	db.ExecInsert(idchan, reschan, insertCmd, roleId, boxContent)

}

func updateHouse(roleId int64, boxContent string) {

	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE treasurehouse SET treasureboxs=? WHERE useruid=?"

	db.ExecUpdate(idchan1, reschan1, updateCmd,
		boxContent, roleId)

}

//判断一个玩家身上曾经是否有宝箱数据
func SaveLogoutBox(boxUid, roleId int64) bool {

	row, ret := db.QueryNormal("SELECT treasureboxs FROM treasurehouse WHERE useruid = ?", roleId)
	if ret == false {
		fmt.Println("从数据库中取出玩家宝箱数据错误")
		return false
	}

	defer row.Close()

	if row == nil || row.Next() == false {
		//在数据库中新添宝箱数据
		itemStr := strconv.Itoa(int(boxUid)) + ","
		insertHouse(roleId, itemStr)
	} else {
		var boxsContent string = ""
		err := row.Scan(&boxsContent)
		if err != nil {
			fmt.Println("解析玩家宝箱数据错误", err)
			return false
		}

		itemStr := strconv.Itoa(int(boxUid)) + ","
		boxsContent += itemStr
		updateHouse(roleId, boxsContent)
	}

	return true

}
