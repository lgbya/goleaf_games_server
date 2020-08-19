package conf

import (
	"encoding/json"
	"github.com/name5566/leaf/log"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dir := filepath.Dir(path)
	dir = "bin/"
	data, err := ioutil.ReadFile(dir+"/conf/server.json")

	if err != nil {
		log.Fatal("%v", err)
	}

	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
	//log.Release("===========解析配置成功=========== \n")
}

