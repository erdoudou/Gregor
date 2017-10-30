package item

import (
	json "encoding/json"
	"fmt"
	io "io/ioutil"
)

func Load(filename string, v interface{}) {
	data, err := io.ReadFile(filename)
	if nil != err {
		fmt.Println(err)
		return
	}

	datajson := []byte(data)
	err = json.Unmarshal(datajson, v)
	if err != nil {
		fmt.Println(err)
		return
	}
}
