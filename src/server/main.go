package main

import (
	"fmt"
	"github.com/name5566/leaf"
	"server/conf"
	"server/game"
	"server/gate"
	"server/login"
	"sync"
	lconf "github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/chanrpc"
)

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)

	//Example()
}

func Example() {
	s := chanrpc.NewServer(10)

	var wg sync.WaitGroup
	wg.Add(1)

	// goroutine 1
	go func() {
		s.Register("f0", func(args []interface{}) {
			fmt.Println("f0")
		})


		wg.Done()

		for {
			s.Exec(<-s.ChanCall)
		}
	}()

	wg.Wait()
	wg.Add(1)

	// goroutine 2
	go func() {
		c := s.Open(10)

		// sync
		c.AsynCall("f0", func(err error) {
			if err != nil {
				fmt.Println(err)
			}
		})
		wg.Done()

	}()

	wg.Wait()

	// Output:
	// 1
	// 1 2 3
	// 3
	// 1
	// 1 2 3
	// 3
}
