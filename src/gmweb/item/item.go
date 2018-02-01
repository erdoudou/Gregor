package item

import (
	"Public/db"
	"com/ViKing/FPSDemo/hallServer/core"
	//"database/sql"
	"fmt"
	//"math/rand"
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
	IsNew      bool
	Paintid    int64
	Badgeid    int64
	Attrs      *Attr
	BaseAttrs  *BaseAttr
}

//数据库查找某个物品dou071025
func FindDbItem(itemuid int64) *Item {
	selectcmd := "SELECT * FROM item WHERE bdel = 0 AND itemuid = ?"
	row, ret := db.QueryNormal(selectcmd, itemuid)

	if ret == false || row == nil {
		fmt.Println("从数据库中取出某个物品数据错误", itemuid)
		return nil
	}

	defer row.Close()

	var itemtempid, bdel, isnew int
	var paintid, badgeid int64

	var hpvalue, bulletloadnumber, movespeed, flexibility, stability, damage, shootspeed, attr8,
		attr9, attr10 float32
	var basehpvalue, basebulletloadnumber, basemovespeed, baseflexibility, basestability, basedamage,
		baseshootspeed, attr18, attr19, attr20 float32

	if row.Next() != false {

		err := row.Scan(&itemuid, &itemtempid, &bdel, &hpvalue, &bulletloadnumber, &movespeed, &flexibility,
			&stability, &damage, &shootspeed, &attr8, &attr9, &attr10, &basehpvalue, &basebulletloadnumber,
			&basemovespeed, &baseflexibility, &basestability, &basedamage, &baseshootspeed, &attr18, &attr19,
			&attr20, &isnew, &paintid, &badgeid)

		if err != nil {
			fmt.Println("解析该物品数据错误", err)
			return nil
		}
		attr := &Attr{
			HpValue:          hpvalue,
			BulletLoadNumber: bulletloadnumber,
			MoveSpeed:        movespeed,
			Flexibility:      flexibility,
			Stability:        stability,
			Damage:           damage,
			ShootSpeed:       shootspeed,
		}

		baseAttr := &BaseAttr{
			BaseHpValue:          basehpvalue,
			BaseBulletLoadNumber: basebulletloadnumber,
			BaseMoveSpeed:        basemovespeed,
			BaseFlexibility:      baseflexibility,
			BaseStability:        basestability,
			BaseDamage:           basedamage,
			BaseShootSpeed:       baseshootspeed,
		}

		prop := &Item{
			itemuid:    itemuid,
			itemtempid: itemtempid,
			bdel:       0,
			Paintid:    paintid,
			Badgeid:    badgeid,
			Attrs:      attr,
			BaseAttrs:  baseAttr,
		}

		if isnew == 0 {
			prop.IsNew = false
		} else {
			prop.IsNew = true
		}

		return prop
	}
	return nil
}

//dou071025
func (prop *Item) GetPropUid() int64 {
	return prop.itemuid

}

//dou071025
func (prop *Item) GetPropTempid() int {
	return prop.itemtempid

}

//由配置文件生成itemdou071025
func CreatAddProp(itemuid int64, itemtempId int) (bool, *Item) {

	ret, propCfg := GetPropCfgByTempID(itemtempId)
	if ret == false {
		fmt.Println("is not find this prop")
		return false, nil
	}
	//进正30%到-30%之间进行随机
	attr := &Attr{
		HpValue:          (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.HpValue,
		BulletLoadNumber: (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.BulletLoadNumber,
		MoveSpeed:        (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.MoveSpeed,
		Flexibility:      (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.Flexibility,
		Stability:        (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.Stability,
		Damage:           (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.Damage,
		ShootSpeed:       (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.ShootSpeed,
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
		IsNew:      true,
	}

	//在数据库中插入一条数据
	item.insertItem()

	return true, item
}

//删除物品dou071026
func (item *Item) DelItem() {
	item.bdel = 1
	item.updateItemBdel()
}

//数据库中生成itemdou071026
func (item *Item) insertItem() {
	idchan := make(chan int, 2)
	reschan := make(chan bool, 2)

	insertCmd := "INSERT  item(itemuid,itemtempid,hpvalue,bulletloadnumber,movespeed,flexibility,stability,damage,shootspeed,basehpvalue,basebulletloadnumber,basemovespeed,baseflexibility,basestability,basedamage,baseshootspeed,isnew) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	db.ExecInsert(idchan, reschan, insertCmd, item.itemuid, item.itemtempid, item.Attrs.HpValue, item.Attrs.BulletLoadNumber,
		item.Attrs.MoveSpeed, item.Attrs.Flexibility, item.Attrs.Stability, item.Attrs.Damage,
		item.Attrs.ShootSpeed, item.Attrs.HpValue, item.Attrs.BulletLoadNumber, item.Attrs.MoveSpeed,
		item.Attrs.Flexibility, item.Attrs.Stability, item.Attrs.Damage, item.Attrs.ShootSpeed, item.IsNew)
}

//数据库更新itemdou071026
func (item *Item) updateItem() {
	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE item SET itemtempid=?,bdel=?,hpvalue=?,bulletloadnumber=?,movespeed=?,flexibility=?,stability=?,damage=?,shootspeed=?,basehpvalue=?,basebulletloadnumber=?,basemovespeed=?,baseflexibility=?,basestability=?,basedamage=?,baseshootspeed=?,isnew=? WHERE itemuid=?"

	db.ExecUpdate(idchan1, reschan1, updateCmd, item.itemtempid, item.bdel, item.Attrs.HpValue,
		item.Attrs.BulletLoadNumber, item.Attrs.MoveSpeed, item.Attrs.Flexibility, item.Attrs.Stability,
		item.Attrs.Damage, item.Attrs.ShootSpeed, item.BaseAttrs.BaseHpValue, item.BaseAttrs.BaseBulletLoadNumber,
		item.BaseAttrs.BaseMoveSpeed, item.BaseAttrs.BaseFlexibility, item.BaseAttrs.BaseStability, item.BaseAttrs.BaseDamage,
		item.BaseAttrs.BaseShootSpeed, item.IsNew, item.itemuid)
}

//改变商品使用状态dou071026
func (item *Item) updateItemBdel() {
	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE item SET bdel=? WHERE itemuid=?"

	db.ExecUpdate(idchan1, reschan1, updateCmd,
		item.bdel, item.itemuid)
}

//物品升级dou071026
func (item *Item) ItemLevelUp() bool {
	fmt.Println("升级物品之前物品的属性", item.Attrs, item.BaseAttrs)
	flag := false
	if item.Attrs.BulletLoadNumber < 2*item.BaseAttrs.BaseBulletLoadNumber {
		flag = true
		num := float32(algorithmMgr.RandInt(5, 10)) / 100.00
		fmt.Println("随机数1：", num)
		if item.Attrs.BulletLoadNumber*num > item.BaseAttrs.BaseBulletLoadNumber*2 {
			item.Attrs.BulletLoadNumber *= 2
		} else {
			item.Attrs.BulletLoadNumber *= (1 + num)
		}
	}

	if item.Attrs.Damage < 2*item.BaseAttrs.BaseDamage {
		flag = true
		num := float32(algorithmMgr.RandInt(5, 10)) / 100.00
		fmt.Println("随机数2：", num)
		if item.Attrs.Damage*num > item.BaseAttrs.BaseDamage*2 {
			item.Attrs.Damage *= 2
		} else {
			item.Attrs.Damage *= (1 + num)
		}

	}
	if item.Attrs.Flexibility < 2*item.BaseAttrs.BaseFlexibility {
		flag = true
		num := float32(algorithmMgr.RandInt(5, 10)) / 100.00
		fmt.Println("随机数3：", num)
		if item.Attrs.Flexibility*num > item.BaseAttrs.BaseFlexibility*2 {
			item.Attrs.Flexibility *= 2
		} else {
			item.Attrs.Flexibility *= (1 + num)
		}

	}
	if item.Attrs.HpValue < 2*item.BaseAttrs.BaseHpValue {
		flag = true
		num := float32(algorithmMgr.RandInt(5, 10)) / 100.00
		fmt.Println("随机数4：", num)
		if item.Attrs.HpValue*num > item.BaseAttrs.BaseHpValue*2 {
			item.Attrs.HpValue *= 2
		} else {
			item.Attrs.HpValue *= (1 + num)
		}

	}

	if item.Attrs.MoveSpeed < 2*item.BaseAttrs.BaseMoveSpeed {
		flag = true
		num := float32(algorithmMgr.RandInt(5, 10)) / 100.00
		fmt.Println("随机数5：", num)
		if item.Attrs.MoveSpeed*num > item.BaseAttrs.BaseMoveSpeed*2 {
			item.Attrs.MoveSpeed *= 2
		} else {
			item.Attrs.MoveSpeed *= (1 + num)
		}

	}

	if item.Attrs.ShootSpeed < 2*item.BaseAttrs.BaseShootSpeed {
		flag = true
		num := float32(algorithmMgr.RandInt(5, 10)) / 100.00
		fmt.Println("随机数6：", num)
		if item.Attrs.ShootSpeed*num > item.BaseAttrs.BaseShootSpeed*2 {
			item.Attrs.ShootSpeed *= 2
		} else {
			item.Attrs.ShootSpeed *= (1 + num)
		}

	}

	if item.Attrs.Stability < 2*item.BaseAttrs.BaseStability {
		flag = true
		num := float32(algorithmMgr.RandInt(5, 10)) / 100.00
		fmt.Println("随机数7：", num)
		if item.Attrs.Stability*num > item.BaseAttrs.BaseStability*2 {
			item.Attrs.Stability *= 2
		} else {
			item.Attrs.Stability *= (1 + num)
		}

	}

	if flag == true {
		fmt.Println("进行玩家数据库数据的更新")
		item.updateAttr()
	}

	return flag
}

//改变物品attr属性dou071026
func (prop *Item) updateAttr() {
	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE item SET hpvalue=?,bulletloadnumber=?,movespeed=?,flexibility=?,stability=?,damage=?,shootspeed=? WHERE itemuid=?"

	db.ExecUpdate(idchan1, reschan1, updateCmd,
		prop.Attrs.HpValue, prop.Attrs.BulletLoadNumber, prop.Attrs.MoveSpeed, prop.Attrs.Flexibility,
		prop.Attrs.Stability, prop.Attrs.Damage, prop.Attrs.ShootSpeed, prop.itemuid)

}

//保存物品数据dou071025
func (item *Item) SaveItemData() {
	row, ret := db.QueryNormal("SELECT * FROM item WHERE itemuid = ?", item.itemuid)
	if ret == false {
		fmt.Println("从数据库中查找物品数据失败")
		return
	}

	defer row.Close()

	if row == nil || row.Next() == false {
		item.insertItem()
	} else {
		item.updateItem()
	}

}

//贴花
func (item *Item) LookNewItem() {
	item.IsNew = false
}

//离线存itemdou071026
func SaveLogoutItem(itemuid int64, itemtempId int) bool {
	//根据配置表得到属性值
	ret, propCfg := GetPropCfgByTempID(itemtempId)
	if ret == false {
		fmt.Println("is not find this prop")
		return false
	}

	//有一个随机公式，将基础属性取出来，进行随机，然后将值存进正30%到-30%之间进行随机

	attr := &Attr{
		HpValue:          (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.HpValue,
		BulletLoadNumber: (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.BulletLoadNumber,
		MoveSpeed:        (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.MoveSpeed,
		Flexibility:      (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.Flexibility,
		Stability:        (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.Stability,
		Damage:           (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.Damage,
		ShootSpeed:       (1 + algorithmMgr.RandFloat(-30, 30)) * propCfg.ShootSpeed,
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
		IsNew:      true,
	}

	//在数据库中插入一条数据
	item.insertItem()

	return true
}
