package item

import (
	"Public/db"
	"fmt"
	"math/rand"
	"strconv"
)

//存储一个物品当前属性
type Attr struct {
	HpValue          float32
	BulletLoadNumber float32
	MoveSpeed        float32
	Flexibility      float32
	Stability        float32
	Damage           float32
	ShootSpeed       float32
}

//存储一个物品的当前基础属性
type BaseAttr struct {
	BaseHpValue          float32
	BaseBulletLoadNumber float32
	BaseMoveSpeed        float32
	BaseFlexibility      float32
	BaseStability        float32
	BaseDamage           float32
	BaseShootSpeed       float32
}

//物品结构
type Item struct {
	itemuid    int64
	itemtempid int
	bdel       int
	Scrawlid   int64
	Paintid    int64
	Badgeid    int64
	Attrs      *Attr
	BaseAttrs  *BaseAttr
}

//根据配置文件生成一个prop
func CreatAddProp(itemuid int64, itemtempId int) (bool, *Item) {
	//根据配置表得到属性值
	ret, propCfg := GetPropCfgByTempID(itemtempId)
	if ret == false {
		fmt.Println("is not find this prop")
		return false, nil
	}

	//有一个随机公式，将基础属性取出来，进行随机，然后将值存进正30%到-30%之间进行随机

	attr := &Attr{
		HpValue:          float32(1-rand.Intn(3)/100) * propCfg.HpValue,
		BulletLoadNumber: float32(1-rand.Intn(3)/100) * propCfg.BulletLoadNumber,
		MoveSpeed:        float32(1-rand.Intn(3)/100) * propCfg.MoveSpeed,
		Flexibility:      float32(1-rand.Intn(3)/100) * propCfg.Flexibility,
		Stability:        float32(1-rand.Intn(3)/100) * propCfg.Stability,
		Damage:           float32(1-rand.Intn(3)/100) * propCfg.Damage,
		ShootSpeed:       float32(1-rand.Intn(3)/100) * propCfg.ShootSpeed,
	}

	baseAttr := &BaseAttr{
		BaseHpValue:          attr.HpValue,
		BaseBulletLoadNumber: attr.BulletLoadNumber,
		BaseMoveSpeed:        attr.MoveSpeed,
		BaseFlexibility:      attr.Flexibility,
		BaseStability:        attr.Stability,
		BaseDamage:           attr.Damage,
		BaseShootSpeed:       attr.ShootSpeed,
	}

	//在内存中加入数据
	var item *Item
	item = &Item{
		itemuid:    itemuid,
		itemtempid: propCfg.ItemTempID,
		bdel:       0,
		Attrs:      attr,
		BaseAttrs:  baseAttr,
	}
	//_propMgr[itemuid] = item

	//在数据库中插入一条数据
	item.insertItemToDb()

	return true, item

}

//在数据库中生成一个prop
func (item *Item) insertItemToDb() {
	idchan := make(chan int, 2)
	reschan := make(chan bool, 2)

	insertCmd := "INSERT  item(itemuid,itemtempid,hpvalue,bulletloadnumber,movespeed,flexibility,stability,damage,shootspeed,basehpvalue,basebulletloadnumber,basemovespeed,baseflexibility,basestability,basedamage,baseshootspeed) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	db.ExecInsert(idchan, reschan, insertCmd, item.itemuid, item.itemtempid, item.Attrs.HpValue, item.Attrs.BulletLoadNumber,
		item.Attrs.MoveSpeed, item.Attrs.Flexibility, item.Attrs.Stability, item.Attrs.Damage,
		item.Attrs.ShootSpeed, item.Attrs.HpValue, item.Attrs.BulletLoadNumber, item.Attrs.MoveSpeed,
		item.Attrs.Flexibility, item.Attrs.Stability, item.Attrs.Damage, item.Attrs.ShootSpeed)

}

//在数据库中更新item的所有数据
func (item *Item) updateItemToDb() {
	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE item SET itemtempid=?,bdel=?,hpvalue=?,bulletloadnumber=?,movespeed=?,flexibility=?,stability=?,damage=?,shootspeed=?,basehpvalue=?,basebulletloadnumber=?,basemovespeed=?,baseflexibility=?,basestability=?,basedamage=?,baseshootspeed=? WHERE itemuid=?"

	db.ExecUpdate(idchan1, reschan1, updateCmd, item.itemtempid, item.bdel, item.Attrs.HpValue,
		item.Attrs.BulletLoadNumber, item.Attrs.MoveSpeed, item.Attrs.Flexibility, item.Attrs.Stability,
		item.Attrs.Damage, item.Attrs.ShootSpeed, item.BaseAttrs.BaseHpValue, item.BaseAttrs.BaseBulletLoadNumber,
		item.BaseAttrs.BaseMoveSpeed, item.BaseAttrs.BaseFlexibility, item.BaseAttrs.BaseStability, item.BaseAttrs.BaseDamage,
		item.BaseAttrs.BaseShootSpeed, item.itemuid)
}

//改变某个商品的使用状态
func (item *Item) updateItemBdel() {

	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE item SET bdel=? WHERE itemuid=?"

	db.ExecUpdate(idchan1, reschan1, updateCmd,
		item.bdel, item.itemuid)

}

//改变一个物品的attr属性
func (prop *Item) updateAttrDb() {
	//fmt.Println("打印出去保存数据")

	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE item SET hpvalue=?,bulletloadnumber=?,movespeed=?,flexibility=?,stability=?,damage=?,shootspeed=? WHERE itemuid=?"

	db.ExecUpdate(idchan1, reschan1, updateCmd,
		prop.Attrs.HpValue, prop.Attrs.BulletLoadNumber, prop.Attrs.MoveSpeed, prop.Attrs.Flexibility,
		prop.Attrs.Stability, prop.Attrs.Damage, prop.Attrs.ShootSpeed, prop.itemuid)

}

//保存物品数据
func (item *Item) SaveItemData() {
	//首先根据物品的uid在数据库中查找是否有该数据
	row, ret := db.QueryNormal("SELECT * FROM item WHERE itemuid = ?", item.itemuid)
	if ret == false {
		fmt.Println("从数据库中查找物品数据失败")
		return
	}

	if row == nil || row.Next() == false {
		fmt.Println("item中没有数据")
		//在数据库中新添item数据
		item.insertItemToDb()
	}

	defer row.Close()

	item.updateItemToDb()

}

//离线存item数据
func SaveLogoutItem(itemuid int64, itemtempId int) bool {
	//根据配置表得到属性值
	ret, propCfg := GetPropCfgByTempID(itemtempId)
	if ret == false {
		fmt.Println("is not find this prop")
		return false
	}

	//有一个随机公式，将基础属性取出来，进行随机，然后将值存进正30%到-30%之间进行随机

	attr := &Attr{
		HpValue:          float32(1-rand.Intn(3)/100) * propCfg.HpValue,
		BulletLoadNumber: float32(1-rand.Intn(3)/100) * propCfg.BulletLoadNumber,
		MoveSpeed:        float32(1-rand.Intn(3)/100) * propCfg.MoveSpeed,
		Flexibility:      float32(1-rand.Intn(3)/100) * propCfg.Flexibility,
		Stability:        float32(1-rand.Intn(3)/100) * propCfg.Stability,
		Damage:           float32(1-rand.Intn(3)/100) * propCfg.Damage,
		ShootSpeed:       float32(1-rand.Intn(3)/100) * propCfg.ShootSpeed,
	}

	baseAttr := &BaseAttr{
		BaseHpValue:          attr.HpValue,
		BaseBulletLoadNumber: attr.BulletLoadNumber,
		BaseMoveSpeed:        attr.MoveSpeed,
		BaseFlexibility:      attr.Flexibility,
		BaseStability:        attr.Stability,
		BaseDamage:           attr.Damage,
		BaseShootSpeed:       attr.ShootSpeed,
	}

	//在内存中加入数据
	var item *Item
	item = &Item{
		itemuid:    itemuid,
		itemtempid: propCfg.ItemTempID,
		bdel:       0,
		Attrs:      attr,
		BaseAttrs:  baseAttr,
	}
	//_propMgr[itemuid] = item

	//在数据库中插入一条数据
	item.insertItemToDb()

	return true
}

func SaveLogoutBag(bagUid, roleId int64) bool {
	//先查找有数据就更新没有数据就增加
	row, ret := db.QueryNormal("SELECT itemslots FROM backpack WHERE useruid = ?", roleId)
	if ret == false {
		fmt.Println("从数据库中查找背包数据失败")
		return false
	}

	defer row.Close()

	if row == nil || row.Next() == false {
		fmt.Println("背包中没有数据")
		//在数据库中新添item数据
		itemStr := strconv.Itoa(int(bagUid)) + ","

		idchan := make(chan int, 2)
		reschan := make(chan bool, 2)

		insertCmd := "INSERT  INTO backpack(useruid,itemslots) VALUE (?,?)"
		db.ExecInsert(idchan, reschan, insertCmd, roleId, itemStr)

	} else {
		var bagContent string = ""
		err := row.Scan(&bagContent)
		if err != nil {
			fmt.Println("解析玩家背包数据错误", roleId, err)
			return false
		}

		itemStr := strconv.Itoa(int(bagUid)) + ","
		bagContent += itemStr

		idchan1 := make(chan int, 2)
		reschan1 := make(chan bool, 2)

		updateCmd := "UPDATE backpack SET itemslots=? WHERE useruid=?"

		db.ExecUpdate(idchan1, reschan1, updateCmd,
			bagContent, roleId)
	}

	return true

}
