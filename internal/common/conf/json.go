package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

var LogFlag = log.Ldate | log.Ltime

var Server struct {
	WSAddr      string
	CertFile    string
	KeyFile     string
	TCPAddr     string
	MaxConnNum  int
	ConsolePort int
	ProfilePath string
	Md5Key      string
	DB          DB
	Log         Log
	Gate        Gate
	Robot       Robot
	Cache       Cache
	Skeleton    Skeleton
}

type Skeleton struct {
	GoLen              int
	TimerDispatcherLen int
	AsynCallLen        int
	ChanRPCLen         int
}

type Log struct {
	Level string
	Path  string
}

type Gate struct {
	PendingWriteNum int
	MaxMsgLen       uint32
	HTTPTimeout     time.Duration
	LenMsgLen       int
	LittleEndian    bool
}

type DB struct {
	Username string
	Password string
	Host     string
	Database string
	MaxIdle  int
	MaxOpen  int
}

type Robot struct {
	Num       int
	MatchTime time.Duration
}

type Cache struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

func init() {

	data, err := ioutil.ReadFile("configs/server.json")

	if err != nil {
		log.Fatal("%v", err)
	}

	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
}
