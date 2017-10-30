package box

import (
	"Public/db"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Box struct {
	BoxUid       int64
	BoxTempId    int //宝箱的id
	IsStart      int
	OpenTime     time.Time
	ItemTempIds  string
	NeedTime     int    // 从数据库读取,生成宝箱时确定(秒)
	itemWeight   string // 物品随机权重
	itemPrice    int    // 售价
	ItemDescribe string // 物品描述
	Platformtype int    //平台信息
	Bdel         int    //是否被开启过
}

//将一个随机宝箱插入数据库/有用/
func (box *Box) insertBox() bool {
	idchan := make(chan int, 2)
	reschan := make(chan bool, 2)

	insertCmd := "INSERT  INTO treasurebox(treasureuid,treasuretempid,isstart,opentime,bdel,platformid,needtime) VALUE (?,?,?,?,?,?,?)"
	db.ExecInsert(idchan, reschan, insertCmd, box.BoxUid, box.BoxTempId, box.IsStart, box.OpenTime, box.Bdel, box.Platformtype, box.NeedTime)

	ret := <-reschan
	if ret {
		return true
	} else {
		return false
	}
}

func (box *Box) GetBoxUid() int64 {
	return box.BoxUid

}

func (box *Box) GetBoxTempid() int32 {
	return int32(box.BoxTempId)

}

func (box *Box) GetBoxOpenTime() string {
	return box.OpenTime.Format("2006-01-02 15:04:05")

}

// 随机掉落宝箱内物品(根据权重)
// 返回值:物品模板ID
func (box *Box) WeightRandomDropItem() (bool, int) {
	//bFind, boxCfgInfo := GetBoxCfgByTempID(int(box.GetBoxTempid()))
	//if !bFind {
	//	fmt.Println("没有找到宝箱配置")
	//	return false, 0
	//}

	itemids := strings.Split(box.ItemTempIds, ",")
	itemweights := strings.Split(box.itemWeight, ",")

	if len(itemids) != len(itemweights) {
		fmt.Println("物品种类和权重数量不匹配,箱子模板ID:", box.GetBoxTempid(), len(itemids), len(itemweights), itemids, itemweights)
		return false, 0
	}

	// 字符串数组转换为int数组
	var itemtemArray []int
	for i := 0; i < len(itemweights); i++ {
		if itemweights[i] != "" {
			itemvalue, _ := strconv.Atoi(itemweights[i])
			itemtemArray = append(itemtemArray, itemvalue)
		}
	}

	// 计算随机索引
	weightIndex := CalculateProbability(itemtemArray)
	itemidStr := itemids[weightIndex]
	itemidInt, _ := strconv.Atoi(itemidStr)
	return true, itemidInt
}

// 概率计算
func CalculateProbability(proArray []int) (arrayindex int) {
	var totalWeight int = 0
	for _, value := range proArray {
		totalWeight += value
	}

	randNum := rand.Intn(totalWeight)

	for i := 0; i < len(proArray); i++ {
		if randNum < proArray[i] {
			return i
		} else {
			randNum -= proArray[i]
		}
	}

	return len(proArray) - 1
}

//更新box所有数据
func (box *Box) updateBox() {
	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE treasurebox SET treasuretempid=?,isstart=?,opentime=?,bdel=?,platformid=?,needtime=? WHERE treasureuid=?"

	db.ExecUpdate(idchan1, reschan1, updateCmd,
		box.BoxTempId, box.IsStart, box.OpenTime, box.Bdel, box.Platformtype, box.NeedTime, box.BoxUid)

}

//保存宝箱数据
func (box *Box) SaveBoxDbData() {
	row, ret := db.QueryNormal("SELECT * FROM treasurebox WHERE treasureuid = ?", box.BoxUid)
	if ret == false {
		fmt.Println("从数据库中查找物品数据失败")
		return
	}

	if row == nil || row.Next() == false {
		fmt.Println("item中没有数据")
		//在数据库中新添item数据
		box.insertBox()
	}

	defer row.Close()

	box.updateBox()
}

//玩家离线存box数据/有用/
func SaveLogoutBoxData(boxuid int64, boxTempId int, platformtype int) (bool, *Box) {
	ret, boxCfg := GetBoxCfgByTempID(boxTempId)
	if ret == false {
		return false, nil
	}
	// 计算宝箱开启需要时间
	//openBoxNeedtime := CalculateOpenBoxNeedTime(boxCfg)

	box := &Box{
		BoxUid:       boxuid,
		BoxTempId:    boxCfg.BoxTempID,
		IsStart:      0,
		ItemTempIds:  _boxCfgMgr[boxTempId].ItemTempIds,
		NeedTime:     60 * 5,
		itemWeight:   _boxCfgMgr[boxTempId].ItemWeight,
		itemPrice:    _boxCfgMgr[boxTempId].ItemPrice,
		ItemDescribe: _boxCfgMgr[boxTempId].ItemDescribe,
		Platformtype: platformtype,
		Bdel:         0,
	}

	box.insertBox()

	return true, box
}

func (box *Box) updateNeedTime() {
	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE treasurebox SET needtime=? WHERE treasureuid=?"

	db.ExecUpdate(idchan1, reschan1, updateCmd, box.NeedTime, box.BoxUid)

}

//改变宝箱的needTime
func (box *Box) UseKey() {
	box.NeedTime = 0
	box.updateNeedTime()
}
