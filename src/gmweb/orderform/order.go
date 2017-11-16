package orderMgr

// 生成订单
import (
	"Public/db"
	//"com/ViKing/FPSDemo/hallServer/item"
	"fmt"
	"gmweb/item"
	"strconv"
)

type orderManager struct {
}

type OrderInfo struct {
	useruid          int64   // 角色GUID
	steamID          string  // steam账号ID
	appID            string  // 应用程序ID
	language         string  // 语言类型
	cointype         string  // 币种
	goodsitemID      int     // 商品模板ID
	goodsitemNum     int     // 商品数量
	totalprice       float32 // 总金额
	goodsinstruction string  // 物品描述
	goodstype        string  // 物品类型
}

// 添加新订单
// useruid int64:服务器uid
func AddNewOrder(useruid int64, steamuseruid string, appid string, itemid int, itemnum int) {

	bSuc, itemtotalPrice := CalculateItemTotalPrice(itemid, itemnum)
	if bSuc == false {
		fmt.Println("计算物品总价出错,终止本次购买")
		return
	}
	// 放进数据库
	var orderMsg OrderInfo
	orderMsg.useruid = useruid
	orderMsg.appID = appid
	orderMsg.steamID = steamuseruid
	orderMsg.goodsitemID = itemid
	orderMsg.goodsitemNum = itemnum
	orderMsg.totalprice = itemtotalPrice

	// 保存到数据库
	orderuid := SaveOrderIntoDB(orderMsg)

	amountStr := strconv.FormatFloat(float64(itemtotalPrice), 'f', 5, 32)
	fmt.Println("总价格为:", amountStr)
	// 向steam平台请求购买
	go TryToCreatShopcarOnSteam(steamuseruid, orderuid, itemid, useruid, amountStr)
}

// 购买成功,发放物品
func (order *orderManager) BuySuccess(orderid int64) {

}

// 购买失败,打印日志
func (order *orderManager) BuyFailed(orderid int64) {

}

// 保存到数据库
func SaveOrderIntoDB(itemOrder OrderInfo) int64 {
	var orderUid int64 = 0
	idchan := make(chan int, 2)
	reschan := make(chan bool, 2)

	insertCmd := `INSERT INTO goodsorder 
	(useruid,steamID,appID,lan,cointype,goodsitemID,goodsitemNum,totalprice) 
 	VALUES (?,?,?,?,?,?,?,?)`
	db.ExecInsert(idchan, reschan, insertCmd, itemOrder.useruid, itemOrder.steamID, itemOrder.appID, itemOrder.language, itemOrder.cointype, itemOrder.goodsitemID, itemOrder.goodsitemNum, itemOrder.totalprice)

	ret := <-reschan
	if ret {

		orderUid = int64(<-idchan)
	}

	return orderUid
}

// 计算物品总价格
func CalculateItemTotalPrice(itemtemid, itemNum int) (bool, float32) {
	ret, itemprice := item.GetItemPriceByTempID(itemtemid)
	if ret == false {
		fmt.Println("计算物品价格时出错,找不到对应配置:", itemtemid)
		return false, 0
	}

	totalPrice := itemprice * float32(itemNum)
	fmt.Println("物品总价为:", totalPrice)
	return true, totalPrice
}
