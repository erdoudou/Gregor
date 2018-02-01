package overallsituation

import (
	"fmt"
	"gmweb/item"
)

var OverallData OverallSituation

type OverallSituation struct {
	MaxBoxSlots       int
	SciencedotNum     int
	MaxBagSlots       int
	LevelUpSciencedot int

	BattlePeriod     int     // 战役周期 单位:小时
	StarUnlockTime   int     // 星系解锁时间 单位:小时
	MaxCommunityAdd  float64 // 最大社区影响力参数1
	MaxCommunityPara float64 // 最大社区影响力参数2
}

func LoadOverallData() {
	item.Load("config/overallsituation.txt", &OverallData)
	fmt.Println("全局配置:", OverallData.MaxBagSlots, OverallData)
}
