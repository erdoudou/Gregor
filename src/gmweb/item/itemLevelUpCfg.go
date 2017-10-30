package item

//import "fmt"

var _propLevelUpCfgMgr map[int]PropLevelUpCfg = make(map[int]PropLevelUpCfg) // 物品配置表

type PropLevelUpCfg struct {
	ItemTempID int // 物品模板ID
	Percent    float64
	Price      int
	Attrs      string
}

func LoadItemLevelUp(filename string) bool {
	var itemData []PropLevelUpCfg

	Load(filename, &itemData)

	for _, item := range itemData {
		_propLevelUpCfgMgr[item.ItemTempID] = item
	}

	//fmt.Println("打印出升级配置文件", _propLevelUpCfgMgr)

	return true
}

// 根据模板ID查找物品配置信息
func GetPropLevelCfgByTempID(itemtempid int) (rst bool, proplevelUpCfg *PropLevelUpCfg) {
	var pItemcfg PropLevelUpCfg
	if _, ok := _propLevelUpCfgMgr[itemtempid]; ok {
		pItemcfg = _propLevelUpCfgMgr[itemtempid]
		return true, &pItemcfg
	}
	return false, nil
}
