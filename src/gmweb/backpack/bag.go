package bagMgr

import (
	"Public/db"
	"database/sql"
	"fmt"
	"gmweb/item"
	"gmweb/overallsituation"
	"strconv"
	"strings"
)

type Bag struct {
	itemContent string // 物品内容集合(逗号分隔符,分隔)e.g:10001,10002, 以逗号结尾
	itemCurNum  int    // 当前物品数量
	itemMgr     map[int64]*item.Item
}

//数据库的背包,物品数据加进内存dou071025
func (bag *Bag) LoadContent(roleid int64) {

	row1 := make(chan *sql.Rows, 1)
	db.Query(row1, "SELECT itemslots FROM backpack WHERE useruid = ?",
		roleid)
	in1 := <-row1

	if in1 == nil {
		return
	}

	defer in1.Close()

	var itemslots string = ""
	if in1.Next() != false {

		err := in1.Scan(&itemslots)

		if err != nil {
			return
		}
	}
	bag.Init(itemslots)
}

//初始化背包dou071025
func (bag *Bag) Init(content string) {
	bag.itemContent = content
	bag.calculateItemNum()
	bag.itemMgr = make(map[int64]*item.Item)
	bag.InitItemMgr()
}

//背包中物品数组赋值dou071025
func (bag *Bag) InitItemMgr() {
	if bag.IsEmpty() == true {
		fmt.Println("查询背包的时候，背包为空")
		return
	}

	for _, v := range strings.Split(bag.itemContent, ",") {
		if v != "" {
			id, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("背包数据的转换错误", err)
				return
			}
			item := item.FindDbItem(int64(id))

			if item != nil {
				bag.itemMgr[int64(id)] = item
			}

		}
	}
}

//计算背包物品数量dou071025
func (bag *Bag) calculateItemNum() {
	if bag.itemContent == "" {
		bag.itemCurNum = 0
		return
	}
	bag.itemCurNum = strings.Count(bag.itemContent, ",")
}

//背包是否为空dou071025
func (bag *Bag) IsEmpty() bool {
	if bag.itemCurNum <= 0 {
		return true
	}

	return false

}

//背包是否有空槽dou071025
func (bag *Bag) HasSlots() bool {
	if overallsituation.OverallData.MaxBagSlots-bag.itemCurNum <= 0 {
		return false
	}

	return true
}

//内存删除物品dou071025
func (bag *Bag) DelItemBagOnly(itemUid int64) {
	itemStr := strconv.Itoa(int(itemUid)) + ","
	bag.itemContent = strings.Replace(bag.itemContent, itemStr, "", -1)
	bag.itemCurNum = bag.itemCurNum - 1

	delete(bag.itemMgr, itemUid)
}

//玩家是否有某个物品dou071025
func (bag *Bag) hasItem(uid int64) bool {
	uidStr := strconv.Itoa(int(uid))
	ret := strings.Contains(bag.itemContent, uidStr)
	if ret {
		return true
	}
	return false
}

//增加一个已经存在的物品(交易)
func (bag *Bag) AddItem(itemuid string) {
	itemUidTrans, _ := strconv.ParseInt(itemuid, 10, 64)
	item := item.FindDbItem(itemUidTrans)

	if item != nil {
		bag.SetItem(item)
	} else {
		fmt.Println("交易物品后添加物品失败,没有找到物品")
	}
}

//增加一个物品dou071025
func (bag *Bag) SetItem(item *item.Item) { //, roleId int64
	itemuid := item.GetPropUid()

	bag.itemMgr[itemuid] = item

	itemStr := strconv.Itoa(int(itemuid)) + ","
	bag.itemContent += itemStr

	bag.itemCurNum = bag.itemCurNum + 1
	//bag.saveBagData(roleId)

}

//背包中删除物品dou071025
func (bag *Bag) DelItem(itemId int64) {
	if bag.IsEmpty() == true || bag.hasItem(itemId) == false {
		return
	}

	if _, ok := bag.itemMgr[itemId]; ok == false {
		return
	}

	item := bag.itemMgr[itemId]
	item.DelItem()

	bag.DelItemBagOnly(itemId)

	//updateBag(roleId, bag.itemContent)
}

//item所有物品数据存储dou071025
func (bag *Bag) SaveAllItemData() {
	if bag.IsEmpty() == true {
		fmt.Println("背包数据为空,玩家没有背包或者物品数据存储")
		return
	}

	for _, v := range bag.itemMgr {
		if v != nil {
			v.SaveItemData()
		}
	}
}

//背包所有物品dou071025
func (bag *Bag) FindAllItem() []*item.Item {
	var itemSli []*item.Item
	if bag.IsEmpty() == true {
		fmt.Println("查询背包的时候，背包为空")
		return itemSli
	}

	for _, v := range bag.itemMgr {
		if v != nil {
			itemSli = append(itemSli, v)
		}
	}

	return itemSli
}

//bag中的某个物品dou071025
func (bag *Bag) FindItem(itemuid int64) *item.Item {

	if bag.IsEmpty() == true || bag.hasItem(itemuid) == false {
		fmt.Println("背包为空或者背包中没有该数据")
		return nil
	}

	item := bag.itemMgr[itemuid]

	return item
}

//保存背包数据dou071025
func (bag *Bag) saveBagData(roleId int64) {
	row, ret := db.QueryNormal("SELECT itemslots FROM backpack WHERE useruid = ?", roleId)
	if ret == false {
		fmt.Println("从数据库中查找背包数据失败")
		return
	}
	defer row.Close()

	if row == nil || row.Next() == false {
		insertBag(roleId, bag.itemContent)
	} else {
		updateBag(roleId, bag.itemContent)
	}

}

//保存背包,item数据dou071025
func (bag *Bag) SaveBagAndItem(roleId int64) {
	if bag.IsEmpty() == true {
		fmt.Println("背包数据为空")
		return
	}

	bag.saveBagData(roleId)
	bag.SaveAllItemData()
}

// 删除物品(参数为字符串,只删除背包数据,不删除物品数据,慎用!)
func (bag *Bag) delItemStr(itemUid string) bool {
	itemUidTrans, _ := strconv.ParseInt(itemUid, 10, 64)
	bag.DelItemBagOnly(itemUidTrans)
	return true
}

// 检测背包过期数据
func (bag *Bag) CheckBackOverdue(steamitem []string) {

	bagitemArray := strings.Split(bag.itemContent, ",")

	// 检测非法物品
	for _, v := range bagitemArray {
		if v != "" {
			if IsExistItem(v, steamitem) == false {
				// 没有在steam上面的物品都是非法物品
				bag.delItemStr(v)
			}
		}
	}

	// 验证遗漏物品
	for _, steamValue := range steamitem {
		if steamValue != "" {
			if IsExistItem(steamValue, bagitemArray) == false {
				// steam 上面有,但在服务器没有的,要添加物品
				bag.AddItem(steamValue)
			}
		}
	}
}

// 检测指定物品是否存在与指定数组中
func IsExistItem(bagitemid string, steamitem []string) bool {
	for _, itemid := range steamitem {
		if strings.Compare(bagitemid, itemid) == 0 {
			return true
		}
	}

	return false
}

//新增背包数据dou071025
func insertBag(roleId int64, itemContent string) {
	idchan := make(chan int, 2)
	reschan := make(chan bool, 2)

	insertCmd := "INSERT  INTO backpack(useruid,itemslots) VALUE (?,?)"
	db.ExecInsert(idchan, reschan, insertCmd, roleId, itemContent)
}

//更新背包数据dou071025
func updateBag(roleId int64, itemContent string) {
	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE backpack SET itemslots=? WHERE useruid=?"

	db.ExecUpdate(idchan1, reschan1, updateCmd,
		itemContent, roleId)
}

//离线存数据dou071025
func SaveLogoutBag(bagUid, roleId int64) bool {
	row, ret := db.QueryNormal("SELECT itemslots FROM backpack WHERE useruid = ?", roleId)
	if ret == false {
		fmt.Println("从数据库中查找背包数据失败")
		return false
	}

	defer row.Close()

	if row == nil || row.Next() == false {
		itemStr := strconv.Itoa(int(bagUid)) + ","
		insertBag(roleId, itemStr)

	} else {
		var bagContent string = ""
		err := row.Scan(&bagContent)
		if err != nil {
			fmt.Println("解析玩家背包数据错误", err)
			return false
		}

		itemStr := strconv.Itoa(int(bagUid)) + ","
		bagContent += itemStr

		updateBag(roleId, bagContent)
	}

	return true
}
