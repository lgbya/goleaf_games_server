package robot

import (
	"github.com/name5566/leaf/network"
	"net"
)

//机器人agent
type robotAgent struct {
	robotCallCh 	 chan interface{}
	conn     network.Conn
	userData interface{}
}

func (a *robotAgent) Run() {

}

func (a *robotAgent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *robotAgent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *robotAgent) Close() {
}

func (a *robotAgent) Destroy() {
}

func (a *robotAgent) WriteMsg(msg interface{}) {
	a.robotCallCh<-msg
}

func (a *robotAgent) UserData() interface{} {
	return a.userData
}

func (a *robotAgent) SetUserData(data interface{}) {
	a.userData = data
}