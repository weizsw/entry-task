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
	"getprofile":     models.GetUserProfile,
	"changenickname": models.ChangeNickname,
	"updateprofile":  models.UpdatePic,
}

func HandleConn(conn net.Conn) {
	for {
		data := make([]byte, 4096)
		session := rpc.NewSession(conn)
		data, err := session.Read()
		if err != nil {
			log.Fatal(err.Error())
		}
		req := &pb.Msg{}
		err = req.Decode(data)

		// log.Println("receive func name:", ret.ServiceName)
		HandlerFunc[req.ServiceName](session, req)
	}
}
