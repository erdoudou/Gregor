package orderMgr

// 付账验证
import (
	"Public/db"
	//"com/ViKing/FPSDemo/hallServer/daily"
	json "encoding/json"
	//	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// 测试地址
	payfinalizeurl = "https://api.steampowered.com/ISteamMicroTxnSandbox/FinalizeTxn/v2?" //"https://partner.steam-api.com/ISteamMicroTxn/FinalizeTxn/v2?"
)

type PaymentResponse struct {
	Response PaymentContent
}

type PaymentContent struct {
	Result string       // 操作结果
	Params PaymentParam // 参数
	Error  PaymentError // 错误信息
}

type PaymentParam struct {
	Orderid string // 订单ID
	Transid string
}

// 支付验证错误信息
type PaymentError struct {
	Errorcode int    // 错误代码
	Errordesc string // 错误描述
}

// 处于验证结果
func PaymentValidate(orderid int64, useruid int64, itemType int) {
	go TryToFinalizePayResult(orderid, useruid, itemType)
}

// 尝试验证付款结果
func TryToFinalizePayResult(orderid int64, useruid int64, itemType int) bool {
	for i := 1; i <= 10; i++ {
		payResult := FinalizePayResultOnsteam(orderid)
		if payResult == ITEM_HTTP_TRADE_ERROR_NULL {
			fmt.Println("玩家已经成功付款,现在授予物品")
			GrantItemToPlayer(orderid, useruid, itemType)
			return true
		} else if payResult != ITEM_HTTP_TRADE_ERROR_NET {
			fmt.Println("验证付款非网络错误,不进行重试")
			return false
		}
		fmt.Println("尝试验证付款结果失败一次,正在进行第", i, "次尝试")
	}
	fmt.Println("多次尝试验证付款结果都失败了")
	// 发送验证失败消息
	return false
}

// 在steam验证订单付款结果
func FinalizePayResultOnsteam(orderid int64) int {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(25 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	form := url.Values{}
	form.Set("key", "DCD9C36F1F54A96F707DFBE833600167")
	form.Set("orderid", strconv.Itoa(int(orderid)))
	form.Set("appid", "511600")
	//form.Set("steamid", steamuserid)

	b := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", payfinalizeurl, b)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("form", "xml")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		fmt.Println("network error!")
		//steamItemDaily.AddRecordeDaily(strconv.FormatInt(orderid, 10), "验证付款结果网络错误", orderid, 0, "id是订单ID")
		return ITEM_HTTP_TRADE_ERROR_NET
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println("验证订单付款结果返回内容:", string(body))
	// add by hjd 20170913
	var payment PaymentResponse
	err = json.Unmarshal([]byte(body), &payment)
	if err != nil {
		fmt.Println("购物车返回内容解析出错:", err, "可能没有完成付款")
		//steamItemDaily.AddRecordeDaily(strconv.FormatInt(orderid, 10), "购物车返回内容解析出错", orderid, 0, err.Error())
		return ITEM_HTTP_TRADE_ERROR_ANALYZE
	}

	if payment.Response.Result != "OK" {
		fmt.Println("验证付款steam平台返回失败:", payment.Response.Error.Errordesc)
		//steamItemDaily.AddRecordeDaily(strconv.FormatInt(orderid, 10), "验证付款steam平台返回失败", orderid, 0, payment.Response.Error.Errordesc)
		return ITEM_HTTP_TRADE_ERROR_LOGIC
	}
	// add end
	//fmt.Println(res.Header)
	if res.Header["X-Error_message"] != nil {
		fmt.Println("验证订单付款结果时出错:", res.Header["X-Error_message"])
		//steamItemDaily.AddRecordeDaily(strconv.FormatInt(orderid, 10), "验证订单付款结果时出错", orderid, 0, res.Header["X-Error_message"][0])
		return ITEM_HTTP_TRADE_ERROR_HEAD
	}

	return ITEM_HTTP_TRADE_ERROR_NULL
}

// 发放物品给玩家
func GrantItemToPlayer(orderid int64, useruid int64, itemType int) {
	// 去数据库查询需要授予哪些物品
	bSuc, steamuserid, itemdefid := SelectItemInOrder(orderid)
	if bSuc == false {
		fmt.Println("查找订单数据失败:", orderid, useruid)
		return
	}

	// 去steam添加物品
	TryToAddItemOnSteam(steamuserid, strconv.Itoa(itemdefid), 0, ITEM_TRADE_TYPE_PAYMENT)

}

// 查询订单中物品及数量
func SelectItemInOrder(orderid int64) (bool, string, int) {
	row, ret := db.QueryNormal("SELECT steamID,goodsitemID FROM goodsorder WHERE orderID = ?", orderid)
	if ret == false {
		fmt.Println("从数据库中取出所有的订单数据错误")
		return false, "", 0
	}

	if row == nil {
		fmt.Println("订单数据库中没有数据")
		return false, "", 0
	}

	defer row.Close()

	var steamuserid string
	var itemdefid int

	if row.Next() != false {

		err := row.Scan(&steamuserid, &itemdefid) //,

		if err != nil {
			fmt.Println("解析订单数据出现错误")
			return false, "", 0
		}
	}

	return true, steamuserid, itemdefid
}
