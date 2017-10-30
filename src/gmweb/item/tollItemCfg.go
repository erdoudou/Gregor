package item

//import "fmt"

// add by hjd 20170905
// 物品类型定义
const (
	ITEM_TYPE_NULL       = iota
	ITEM_TYPE_WEAPON     // 武器1
	ITEM_TYPE_CONSUME    // 消耗品2
	ITEM_TYPE_KEY        // 钥匙3
	ITEM_TYPE_BOX        // 箱子4
	ITEM_TYPE_SCRAWL     // 涂装5
	ITEM_TYPE_PAINT      // 喷漆6
	ITEM_TYPE_BADGE      // 徽章7
	ITEM_TYPE_SCIENCEDOT // 科研点8
	ITEM_TYPE_MAXNUM
)

// add end

var _propCfgMgr map[int]PropCfg = make(map[int]PropCfg) // 物品配置表

type PropCfg struct {
	ItemTempID int // 物品模板ID
	// add by hjd 20170905
	ItemType      int     // 物品类型
	ItemPrice     float32 // 物品价格
	ItemDescribe  string  // 物品描述
	ItemPlusValue int     // 增益值(e.g:科研点)
	// add end
	HpValue          float32
	BulletLoadNumber float32
	MoveSpeed        float32
	Flexibility      float32
	Stability        float32
	Damage           float32
	ShootSpeed       float32
	Attr8            float32
	Attr9            float32
	Attr10           float32
	IsRepeat         int
}

func LoadItem(filename string) bool {
	var itemData []PropCfg

	//JsonParse := NewJsonStruct()
	Load(filename, &itemData)

	for _, item := range itemData {
		_propCfgMgr[item.ItemTempID] = item
	}

	//fmt.Println(_propCfgMgr)

	return true
}

// 根据模板ID查找物品配置信息
func GetPropCfgByTempID(itemtempid int) (rst bool, propCfg *PropCfg) {
	var pItemcfg PropCfg
	if _, ok := _propCfgMgr[itemtempid]; ok {
		pItemcfg = _propCfgMgr[itemtempid]
		return true, &pItemcfg
	}
	return false, nil
}

func GetPropIsrepeat(itemtempid int) (ret bool, isrepeat int) {
	var pItemcfg PropCfg
	if _, ok := _propCfgMgr[itemtempid]; ok {
		pItemcfg = _propCfgMgr[itemtempid]
		return true, pItemcfg.IsRepeat
	}
	return false, 100
}

// 根据模板id获取物品类型
func GetItemtypeByTempID(itemtempid int) (itemtype int) {
	var pItemcfg PropCfg
	if _, ok := _propCfgMgr[itemtempid]; ok {
		pItemcfg = _propCfgMgr[itemtempid]
		return pItemcfg.ItemType
	}
	return ITEM_TYPE_NULL
}

// 根据模板id获取物品类型
func GetItemPriceByTempID(itemtempid int) (bFind bool, itemprice float32) {
	var pItemcfg PropCfg
	if _, ok := _propCfgMgr[itemtempid]; ok {
		pItemcfg = _propCfgMgr[itemtempid]
		return true, pItemcfg.ItemPrice
	}
	return false, 0
}

// 根据模板id获取物品描述
func GetItemDescByTempID(itemtempid int) (bFind bool, itemdesc string) {
	var pItemcfg PropCfg
	if _, ok := _propCfgMgr[itemtempid]; ok {
		pItemcfg = _propCfgMgr[itemtempid]
		return true, pItemcfg.ItemDescribe
	}
	return false, ""
}

// 根据模板id获取物品增益值
func GetItemPlusByTempID(itemtempid int) (bFind bool, itemvalue int) {
	var pItemcfg PropCfg
	if _, ok := _propCfgMgr[itemtempid]; ok {
		pItemcfg = _propCfgMgr[itemtempid]
		return true, pItemcfg.ItemPlusValue
	}
	return false, 0
}
