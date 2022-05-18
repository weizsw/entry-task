package bootstrap

import (
	"log"
	"net"
	_ "net/http/pprof"

	"github.com/weizsw/entry-task/client/resource"
	"github.com/weizsw/entry-task/client/router"
)

func Init() {
	go initConnPool()
	router.Init()
}

func initConnPool() {
	cp, err := resource.NewConnPool(func() (net.Conn, error) {
		conn, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8989")
		if err != nil {
			return nil, err
		}
		return net.DialTCP("tcp", nil, conn)
	}, 1000)
	if err != nil {
		log.Println(err.Error())
	}

	resource.CP = cp
}
