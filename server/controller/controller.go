package controller

import (
	"log"
	"net"

	"github.com/weizsw/entry-task/pb"
	"github.com/weizsw/entry-task/rpc"
	"github.com/weizsw/entry-task/server/models"
)

var HandlerFunc = map[string]func(*rpc.Session, *pb.Msg){
	"login":          models.Login,
	"register":       models.Register,
	"userinfo":       models.UserInfo,
	"changenickname": models.ChangeNickname,
	"updateprofile":  models.UpdateProfile,
}

func HandleConn(conn net.Conn) {
	for {
		data := make([]byte, 4096)
		session := rpc.NewSession(conn)
		data, err := session.Read()
		if err != nil {
			log.Fatal(err.Error())
		}
		ret := &pb.Msg{}
		err = ret.Decode(data)

		// log.Println("receive func name:", ret.ServiceName)
		HandlerFunc[ret.ServiceName](session, ret)
	}
}
