package gamedata

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/recordfile"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
)

func readRf(st interface{}) *recordfile.RecordFile {
	rf, err := recordfile.New(st)
	if err != nil {
		log.Fatal("%v", err)
	}

	fn := reflect.TypeOf(st).Name() + ".csv"
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dir := filepath.Dir(path)
	//dir = "bin"
	err = rf.Read(dir+"/gamedata/" + fn)
	if err != nil {
		log.Fatal("%v: %v", fn, err)
	}

	return rf
}
