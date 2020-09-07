package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main()  {
	fmt.Println("编译协议")
	file, err := ioutil.ReadFile("api/protocol.json")
	data := make(map[string]interface{})
	err = json.Unmarshal(file, &data)

	if err != nil {
		fmt.Println(err)
	}

	protocol := "package protocol\n\n"

	register := "package protocol\n\n"


	register += "func init() {\n"


	for k, v := range  data{
		vv := v.(map[string]interface{})


		protocol += "type " + k + " struct {\n"
		if vv["field"] != nil {
			fields := vv["field"].(map[string]interface{})
			for k2, v2 := range fields {
				vv2 := v2.(map[string]interface{})
				//typeStr :=
				protocol +=  "	" + k2 + "	" + vv2["type"].(string) + " `json:\""+ vv2["tag"].(string) +"\"" +"`\n"
			}
		}
		protocol += "}\n\n"

		register+= "	Processor.Register(&" + k + "{})\n"

	}
	register += "}"
	_ = ioutil.WriteFile("internal/gate/protocol/register.go", []byte(register), 0666) //写入文件(字节数组)
	_ = ioutil.WriteFile("internal/gate/protocol/protocol.go", []byte(protocol), 0666) //写入文件(字节数组)
	// 关闭文件

}