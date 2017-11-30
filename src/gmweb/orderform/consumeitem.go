package orderMgr

// 消耗物品
import (
	//"com/ViKing/FPSDemo/hallServer/daily"
	json "encoding/json"
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
	consumeurl = "https://api.steampowered.com/IInventoryService/ConsumeItem/v0001?"
)

// 多次尝试消耗物品
func TryToConsumeItemOnSteam(steamuserid string, itemid int64) bool {

	for i := 1; i <= 10; i++ {
		if ConsumeItemIntosteam(steamuserid, itemid) {
			return true
		}
		fmt.Println("失败一次,正在进行第", i, "次尝试")
	}
	fmt.Println("多次尝试删除物品都失败了")
	return false
}

// 在steam消耗物品
func ConsumeItemIntosteam(steamuserid string, itemid int64) bool {
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
	form.Set("appid", "675310")
	form.Set("steamid", steamuserid)
	form.Set("itemid", strconv.FormatInt(itemid, 10))
	form.Set("quantity", "1") // 这里只消耗一个箱子就可以了

	b := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", consumeurl, b)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		fmt.Println("network error!")
		//日志系统暂时屏蔽
		//steamItemDaily.AddRecordeDaily(steamuserid, "消耗物品网络错误", itemid, 0, err.Error())
		return false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(res.Header)

	if res.Header["X-Error_message"] != nil {
		fmt.Println("消耗物品时出错:", res.Header["X-Error_message"])
		//屏蔽日志系统
		//steamItemDaily.AddRecordeDaily(steamuserid, "消耗物品时头部信息出错", itemid, 0, res.Header["X-Error_message"][0])
		return false
	}

	// 解析物品返回
	var gainedItem AddItem_Response
	datajson := []byte(body)
	err = json.Unmarshal(datajson, &gainedItem)
	if err != nil {
		fmt.Println("消耗物品返回解析出错1", err)
		//屏蔽日志系统
		//steamItemDaily.AddRecordeDaily(steamuserid, "消耗物品返回解析出错", itemid, 0, err.Error())
		return false
	}
	fmt.Println(gainedItem.Response.Item_json)

	// steam 是以字符串方式返回的物品定义,导致需要解析两次
	var itemDetail []Item_pro
	datajson = []byte(gainedItem.Response.Item_json)
	err = json.Unmarshal(datajson, &itemDetail)
	if err != nil {
		fmt.Println("消耗物品返回解析出错2", err)
		//屏蔽日志系统
		//steamItemDaily.AddRecordeDaily(steamuserid, "消耗物品返回解析出错", itemid, 0, err.Error())
		return false
	}

	fmt.Println(itemDetail)
	fmt.Println("消耗物品成功")
	return true
}
