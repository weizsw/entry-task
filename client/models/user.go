package models

import (
	"errors"
	"log"
	"net"

	"github.com/weizsw/entry-task/client/errno"
	"github.com/weizsw/entry-task/client/resource"
	"github.com/weizsw/entry-task/pb"
	"github.com/weizsw/entry-task/rpc"
)

func Login(username string, password string) (int, error) {
	conn, err := resource.CP.Get()
	defer resource.CP.Put(conn)

	if err != nil {
		log.Println(err.Error())
		return errno.StatusParamsError, err
	}

	// conn, err := net.DialTimeout("tcp", "127.0.0.1:8989", time.Second)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return errno.StatusServerError, err
	// }
	// defer conn.Close()

	session := rpc.NewSession(conn.(net.Conn))

	msg := pb.Msg{ServiceName: "login", LoginRequest: &pb.LoginRequest{Username: username, Password: password}}

	ret, err := rpc.SendRequestForResponse(session, &msg)

	if ret.Response.Code != 0 {
		return int(ret.Response.Code), err
	}

	return errno.StatusOK, nil
}

func Register(username string, password string) (int, error) {
	conn, err := resource.CP.Get()
	defer resource.CP.Put(conn)

	if err != nil {
		log.Println(err.Error())
		return errno.StatusParamsError, err
	}

	session := rpc.NewSession(conn.(net.Conn))

	msg := pb.Msg{ServiceName: "register", RegisterRequest: &pb.RegisterRequest{Username: username, Password: password}}

	ret, err := rpc.SendRequestForResponse(session, &msg)

	if ret.Response.Code != 0 {
		return int(ret.Response.Code), err
	}

	return errno.StatusOK, nil
}

func GetUserInfo(username string) (int, *pb.Msg, error) {
	conn, err := resource.CP.Get()
	defer resource.CP.Put(conn)

	if err != nil {
		log.Println(err.Error())
		return errno.StatusParamsError, nil, err
	}

	session := rpc.NewSession(conn.(net.Conn))

	msg := pb.Msg{ServiceName: "userinfo", UserInfoRequest: &pb.UserInfoRequest{Username: username}}

	ret, err := rpc.SendRequestForResponse(session, &msg)

	if err != nil {
		return errno.StatusServerError, nil, err
	}

	if ret.Response != nil && ret.Response.Code != 0 {
		return errno.StatusParamsError, nil, errors.New(errno.ErrorMsg[errno.StatusParamsError])
	}

	return errno.StatusOK, ret, nil
}

func ChangeNickname(username string, nickname string) (int, error) {
	conn, err := resource.CP.Get()
	defer resource.CP.Put(conn)

	if err != nil {
		log.Println(err.Error())
		return errno.StatusParamsError, err
	}

	session := rpc.NewSession(conn.(net.Conn))

	msg := pb.Msg{ServiceName: "changenickname", ChangeNicknameRequest: &pb.ChangeNicknameRequest{Username: username, Nickname: nickname}}

	ret, err := rpc.SendRequestForResponse(session, &msg)

	if ret.Response.Code != 0 {
		return int(ret.Response.Code), err
	}

	return errno.StatusOK, nil
}

func UpdateProfile(username string, pic string) (int, error) {
	conn, err := resource.CP.Get()
	defer resource.CP.Put(conn)

	if err != nil {
		log.Println(err.Error())
		return errno.StatusParamsError, err
	}

	session := rpc.NewSession(conn.(net.Conn))

	msg := pb.Msg{ServiceName: "updateprofile", UpdateProfileRequest: &pb.UpdateProfileRequest{Username: username, Pic: pic}}

	ret, err := rpc.SendRequestForResponse(session, &msg)

	if ret.Response.Code != 0 {
		return int(ret.Response.Code), err
	}

	return errno.StatusOK, nil
}
