package box

import (
	"fmt"
)

func LoadCfg() bool {
	bRst := LoadItem("config/box.txt")
	if false == bRst {
		fmt.Println("加载宝箱config数据出错")
		return false
	}

	return true
}
