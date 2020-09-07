package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main()  {
	fmt.Println("编译路由")
	file, err := ioutil.ReadFile("api/protocol.json")
	data := make(map[string]interface{})
	err = json.Unmarshal(file, &data)

	if err != nil {
		fmt.Println(err)
	}

	router := "package gate\n\n"
	router += "import (\n"
	router += "	\"server/internal/game\"\n"
	router += "	\"server/internal/gate/protocol\"\n"
	router += "	\"server/internal/login\"\n"
	router += ")\n\n"

	router += "func init() {\n"

	for k, v := range  data{
		vv := v.(map[string]interface{})



		if vv["router"] != nil {
			router += "	protocol.Processor.SetRouter(&protocol." + k + "{}, " + vv["router"].(string) + ".ChanRPC)\n"
		}

	}
	router += "}"
	_ = ioutil.WriteFile("internal/gate/router.go", []byte(router), 0666) //写入文件(字节数组)

}