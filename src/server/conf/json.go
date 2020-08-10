package conf

import (
	"encoding/json"
	"github.com/name5566/leaf/log"
	"io/ioutil"
)

var Server struct {
	LogLevel    string
	LogPath     string
	WSAddr      string
	CertFile    string
	KeyFile     string
	TCPAddr     string
	MaxConnNum  int
	ConsolePort int
	ProfilePath string
	Md5Key 		string
	SqlSrv		SqlSrv

}

type SqlSrv struct {
	Username 	string
	Password	string
	Host		string
	Database	string
	MaxIdle 	int
	MaxOpen		int
}

func init() {
	data, err := ioutil.ReadFile("bin/conf/server.json")

	if err != nil {
		log.Fatal("%v", err)
	}

	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
	//log.Release("===========解析配置成功=========== \n")
}
