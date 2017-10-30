package item

import (
	"fmt"
)

func LoadCfg() bool {
	bRst := LoadItem("config/item.txt")
	if false == bRst {
		fmt.Println("加载物品config出错")
		return false
	}

	bRst = LoadItemLevelUp("config/levelUp.txt")
	if false == bRst {
		fmt.Println("加载升级config出错")
		return false
	}

	return true
}
