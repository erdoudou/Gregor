package orderMgr

// 购物车
import (
	//	json "encoding/json"
	//"com/ViKing/FPSDemo/hallServer/daily"
	//"com/ViKing/FPSDemo/hallServer/item"
	//"com/ViKing/FPSDemo/hallServer/roleMgr"
	//"com/ViKing/FPSDemo/protobuf/goBin"
	json "encoding/json"
	"fmt"
	"gmweb/item"
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
	shoppingcarurl = "https://api.steampowered.com/ISteamMicroTxnSandbox/InitTxn/V0002?" //"https://api.steampowered.com/ISteamMicroTxn/InitTxn/V0002?"
)

type itemCarResponse struct {
	Response ResponseContent
}

type ResponseContent struct {
	Result string       // 操作结果
	Params itemCarParam // 参数
	Error  itemCarError // 购物车出错原因
}

type itemCarParam struct {
	Orderid string // 订单ID
}

type itemCarError struct {
	Errorcode int    // 错误代码
	Errordesc string // 错误描述
}

// 多次尝试创建购物车
func TryToCreatShopcarOnSteam(steamuserid string, orderid int64, itemid int, useruid int64, amount string) bool {

	for i := 1; i <= 10; i++ {
		createRet := CreatShopcarIntosteam(steamuserid, orderid, itemid, amount)
		if createRet == ITEM_HTTP_TRADE_ERROR_NULL {
			fmt.Println("购物车创建成功,回复客户端")
			//ReplyClientItemCarResult(useruid, orderid, protoMsg.OperateResult_SUCCESE)
			return true
		} else if createRet != ITEM_HTTP_TRADE_ERROR_NET {
			fmt.Println("非网络错误,不进行重试")
			//ReplyClientItemCarResult(useruid, orderid, protoMsg.OperateResult_FAILED)
			return false
		}
		fmt.Println("失败一次,正在进行第", i, "次尝试")
	}
	fmt.Println("多次尝试创建购物车都失败了")
	//ReplyClientItemCarResult(useruid, orderid, protoMsg.OperateResult_FAILED)
	return false
}

// 在steam创建购物车
func CreatShopcarIntosteam(steamuserid string, orderid int64, itemid int, amount string) int {
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
	form.Set("steamid", steamuserid)
	form.Set("appid", "675310")
	form.Set("itemcount", "1")
	form.Set("language", "EN")
	form.Set("currency", "USD")
	form.Set("itemid[0]", strconv.Itoa(itemid))
	form.Set("qty[0]", "1")
	form.Set("amount[0]", amount) //"199"

	bFind, itemDesc := item.GetItemDescByTempID(itemid)
	if bFind == false || itemDesc == "" {
		itemDesc = "this is a mystery item"
	}
	form.Set("description[0]", itemDesc)
	form.Set("category[0]", "Hat")

	b := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", shoppingcarurl, b)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		fmt.Println("network error!")
		//steamItemDaily.AddRecordeDaily(steamuserid, "创建购物车网络错误", orderid, itemid, "")
		return ITEM_HTTP_TRADE_ERROR_NET
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println("购物车返回内容:", string(body))

	// add by hjd 20170913
	var itemResponse itemCarResponse
	err = json.Unmarshal([]byte(body), &itemResponse)
	if err != nil {
		fmt.Println("购物车返回内容解析出错:", err)
		//steamItemDaily.AddRecordeDaily(steamuserid, "购物车返回内容解析出错", orderid, itemid, err.Error())
		return ITEM_HTTP_TRADE_ERROR_ANALYZE
	}
	fmt.Println("返回结果:", itemResponse.Response.Result)
	if itemResponse.Response.Result == "Failure" {
		fmt.Println("错误原因:", itemResponse.Response.Error.Errordesc)
		//steamItemDaily.AddRecordeDaily(steamuserid, "购物车返回返回结果", orderid, itemid, itemResponse.Response.Error.Errordesc)
		return ITEM_HTTP_TRADE_ERROR_LOGIC
	}
	// add end

	if res.Header["X-Error_message"] != nil {
		fmt.Println("创建购物车时出错:", res.Header["X-Error_message"])
		//steamItemDaily.AddRecordeDaily(steamuserid, "创建购物车时头部信息出错", orderid, itemid, res.Header["X-Error_message"][0])
		return ITEM_HTTP_TRADE_ERROR_HEAD
	}

	return ITEM_HTTP_TRADE_ERROR_NULL
}

// 回复客户端创建购物车结果
/*func ReplyClientItemCarResult(useruid int64, orderid int64, result protoMsg.OperateResultRetType) {
	role, hasRole := roleMgr.GetRole(useruid)
	if hasRole == false {
		fmt.Println("角色可能不在线", useruid)
		return
	}

	rspMsg := &protoMsg.ToClientMsg{
		Type: protoMsg.ToClientMsg_PRODUCE_ORDER_RSP,
		ProduceOrderRsp: &protoMsg.ProduceOrderRsp{
			Ret:     result,
			Orderid: orderid,
		},
	}

	role.SendToClient(rspMsg)
}*/
