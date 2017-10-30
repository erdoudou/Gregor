package box

//"fmt"

//import "fmt"

var _boxCfgMgr map[int]BoxCfg = make(map[int]BoxCfg) // 物品配置表

type BoxCfg struct {
	BoxTempID    int // 物品模板ID
	NeedTime     string
	ItemTempIds  string
	ItemWeight   string // 物品随机权重
	ItemPrice    int    // 价格
	ItemDescribe string // 描述
}

func LoadItem(filename string) bool {
	var boxData []BoxCfg

	//JsonParse := NewJsonStruct()
	Load(filename, &boxData)

	for _, box := range boxData {
		_boxCfgMgr[box.BoxTempID] = box
	}

	return true
}

// 根据模板ID查找物品配置信息
func GetBoxCfgByTempID(boxtempid int) (rst bool, boxCfg *BoxCfg) {
	var boxcfg BoxCfg
	if _, ok := _boxCfgMgr[boxtempid]; ok {
		boxcfg = _boxCfgMgr[boxtempid]
		return true, &boxcfg
	}
	return false, nil
}
