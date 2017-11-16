package orderMgr

// 添加物品
import (
	"bytes"
	json "encoding/json"
	"fmt"
	"gmweb/box"
	"gmweb/item"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"Public/db"
	//"com/ViKing/FPSDemo/protobuf/goBin"
)

const (
	//LHST = "http://localhost:8282/IInventoryService/AddItem/v1?key=xxxxxx/"
	//LHST = "http://api.steampowered.com/IInventoryService/AddItem/v1?key=DCD9C36F1F54A96F707DFBE833600167"
	//LHST = "http://api.steampowered.com/IInventoryService/AddItem/v1?"
	LHST = "https://partner.steam-api.com/IInventoryService/AddItem/v1?"
)

// 添加物品类型
const (
	ITEM_TRADE_TYPE_NULL       = iota
	ITEM_TRADE_TYPE_BOX_CHANGE // 箱子兑换
	ITEM_TRADE_TYPE_BOX_OPEN   // 开箱子
	ITEM_TRADE_TYPE_PAYMENT    // 付费购买
	ITEM_TRADE_TYPE_MAX
)

// 添加物品错误类型
const (
	ITEM_HTTP_TRADE_ERROR_NULL    = iota
	ITEM_HTTP_TRADE_ERROR_NET     // 网络错误
	ITEM_HTTP_TRADE_ERROR_ANALYZE // 解析错误
	ITEM_HTTP_TRADE_ERROR_LOGIC   // 逻辑错误
	ITEM_HTTP_TRADE_ERROR_HEAD    // 头部有错误提示
)

// 生成物品结构
type Item_pro struct {
	Accountid               string // 账号ID
	Itemid                  string // 物品UID
	Quantity                int    // 品质(绿,橙,红,传奇)
	Originalitemid          string // 初始物品ID(保留)
	Itemdefid               string // 物品模板ID
	Appid                   int    // appID
	Acquired                string // 获得时间
	State                   string // 当前状态
	Origin                  string // 来源(交换，购买)
	State_changed_timestamp string // 状态变更时间
}

// 回复内容结构体
type Response_Content struct {
	Item_json string //[]Item_pro // 物品数组
}

// 添加物品回复结构体
type AddItem_Response struct {
	Response Response_Content // 回复内容
}

// 多次尝试
// 没有需要删除的物品则delItemID = 0/有用/
func TryToAddItemOnSteam(steamuserid string, itemarray string, delItemID int64, itemTradeType int) bool {

	for i := 1; i <= 10; i++ {
		if AddItemIntosteam(steamuserid, itemarray, delItemID, itemTradeType) {
			return true
		}
		fmt.Println("失败一次,正在进行第", i, "次尝试")
	}
	fmt.Println("多次尝试添加数据都失败了")
	return false
}

// 向steam物品栏添加物品/有用/
func AddItemIntosteam(steamuserid string, itemarray string, delItemID int64, itemTradeType int) bool {

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
	form.Set("appid", "511600")
	form.Set("steamid", steamuserid) //76561198382627854

	/*for index, itemtemid := range itemarray {
		itemEntryStr := GetItemEntry(index)
		fmt.Println(itemEntryStr, "index:", index, "itemtemid:", itemtemid)
		form.Set(itemEntryStr, itemtemid)
	}*/
	form.Set("itemdefid[0]", itemarray)
	b := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", LHST, b)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		//os.Exit(0)
		fmt.Println("network error!")
		//屏蔽日志系统
		//steamItemDaily.AddRecordeDaily(steamuserid, "添加物品网络消息错误", 0, 0, itemarray)
		return false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if res.Header["X-Error_message"] != nil {
		fmt.Println("授予物品时出错:")
		//屏蔽日志系统
		//steamItemDaily.AddRecordeDaily(steamuserid, "授予物品时出错", 0, 0, res.Header["X-Error_message"][0])
		return false
	}

	// 解析物品返回
	var gainedItem AddItem_Response
	datajson := []byte(body)
	err = json.Unmarshal(datajson, &gainedItem)
	if err != nil {
		fmt.Println("解析出错1", err)
		//屏蔽日志系统
		//steamItemDaily.AddRecordeDaily(steamuserid, "授予物品时出错", 0, 0, err.Error())
		return false
	}
	fmt.Println(gainedItem.Response.Item_json)

	// steam 是以字符串方式返回的物品定义,导致需要解析两次
	var itemDetail []Item_pro
	datajson = []byte(gainedItem.Response.Item_json)
	err = json.Unmarshal(datajson, &itemDetail)
	if err != nil {
		fmt.Println("解析出错2", err)
		//屏蔽日志系统
		//steamItemDaily.AddRecordeDaily(steamuserid, "授予物品时出错", 0, 0, err.Error())
		return false
	}

	fmt.Println(itemDetail)

	// 物品进入宝箱容器
	HandlePlatformData(itemDetail, delItemID, itemTradeType)

	return true
}

// 获取物品字段
func GetItemEntry(index int) string {
	temByte := bytes.Buffer{}
	temByte.WriteString("itemdefid[")
	temByte.WriteString(strconv.Itoa(index))
	temByte.WriteString("]")
	return temByte.String()
}

// 处理平台反馈回来的箱子数据/有用/
func HandlePlatformData(itemarray []Item_pro, delItemID int64, itemTradeType int) {

	// 放入宝箱背包
	for _, itemIns := range itemarray {
		itemuid, _ := strconv.ParseInt(itemIns.Itemid, 10, 64)
		itemtempid, _ := strconv.Atoi(itemIns.Itemdefid)
		// 解析操作结果,返回客户端
		roleuid := GetRoleIDbyPlatformID(itemIns.Accountid, int(4)) //4表示steram平台

		fmt.Println("获得物品ID:", itemuid)
		itemType := item.GetItemtypeByTempID(itemtempid)
		fmt.Println("物品类型为:", itemType)

		if itemType == item.ITEM_TYPE_BOX {
			PlayerOffLineBoxDataSave(roleuid, itemuid, itemtempid, int(4)) //4代表protoMsg.PlatformType_PLATFORM_STEAM

		} else {
			PlayerOffLineItemDataSave(roleuid, itemuid, itemtempid)
		}

	}
}

// 获取角色UID
func GetRoleIDbyPlatformID(platformid string, plateType int) int64 {
	//platformidstr := strconv.Itoa(int(platformid))
	row, ret := db.QueryNormal("SELECT globaluid from accountconvert where platformuid = ? and platformtype = ?", platformid, plateType)
	if ret == false {
		fmt.Println("获取角色ID出错")
		return 0
	}

	if row == nil {
		fmt.Println("没有找到对应的角色ID")
		return 0
	}

	defer row.Close()

	var roleuid int64 = 0

	if row.Next() != false {
		err := row.Scan(&roleuid)

		if err != nil {
			fmt.Println("没有找到对对应角色ID")
			return 0
		}
	}
	fmt.Println("查询角色UID:", roleuid)
	return roleuid
}

// 处理玩家不在线时宝箱数据处理
func PlayerOffLineBoxDataSave(roleuid int64, itemuid int64, itemtempid int, plateformType int) {
	bSuc, _ := box.SaveLogoutBoxData(itemuid, itemtempid, plateformType)
	if bSuc {
		fmt.Println("玩家不在线,宝箱直接存入数据库")
	} else {
		//steamItemDaily.AddRecordeDaily(strconv.FormatInt(roleuid, 10), "玩家不在线,宝箱直接存入数据库错误了", itemuid, itemtempid, "id是roleid")
		fmt.Println("玩家不在线,宝箱直接存入数据库错误了")
	}

	bSuc = box.SaveLogoutBox(itemuid, roleuid)
	if bSuc {
		fmt.Println("玩家不在线,宝箱直接存入宝箱背包数据库")
	} else {
		fmt.Println("玩家不在线,宝箱直接存入宝箱背包数据库错误了")
		//steamItemDaily.AddRecordeDaily(strconv.FormatInt(roleuid, 10), "玩家不在线,宝箱直接存入宝箱背包数据库错误了", itemuid, itemtempid, "id是roleid")
	}
}

// 处理玩家不在线时物品数据数据处理
func PlayerOffLineItemDataSave(roleuid int64, itemuid int64, itemtempid int) {
	bSuc := item.SaveLogoutItem(itemuid, itemtempid)

	if bSuc {
		fmt.Println("玩家不在线,物品直接存入数据库")
	} else {
		//steamItemDaily.AddRecordeDaily(strconv.FormatInt(roleuid, 10), "玩家不在线,物品直接存入数据库出错", itemuid, itemtempid, "id是roleid")
		fmt.Println("玩家不在线,物品直接存入数据库出错")
	}

	bSuc = item.SaveLogoutBag(itemuid, roleuid)
	if bSuc {
		fmt.Println("玩家不在线,物品直接存入数据库")
	} else {
		fmt.Println("玩家不在线,物品直接存入背包数据库")
		//steamItemDaily.AddRecordeDaily(strconv.FormatInt(roleuid, 10), "玩家不在线,物品直接存入背包数据库", itemuid, itemtempid, "id是roleid")
	}
}
