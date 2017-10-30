package box

import (
	json "encoding/json"
	"fmt"
	io "io/ioutil"
)

func Load(filename string, v interface{}) {
	data, err := io.ReadFile(filename)
	if nil != err {
		fmt.Println("解析出错了：", err)
		return
	}

	datajson := []byte(data)
	err = json.Unmarshal(datajson, v)
	if err != nil {
		fmt.Println("json出错了", err, filename)
		return
	}
}
