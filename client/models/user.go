package models

import (
	"log"

	"github.com/weizsw/entry-task/client/errno"
	"github.com/weizsw/entry-task/client/resource"
	"github.com/weizsw/entry-task/pb"
	"github.com/weizsw/entry-task/rpc"
)

func Login(username string, password string) (int, interface{}, error) {
	conn, err := resource.CP.Get()
	if err != nil {
		log.Println(err.Error())
		return errno.StatusServerError, nil, err
	}
	defer resource.CP.Put(conn)

	session := rpc.NewSession(conn)

	msg := pb.Msg{ServiceName: "login", LoginRequest: &pb.LoginRequest{Username: username, Password: password}}

	ret, err := rpc.SendRequestForResp(session, &msg)

	if err != nil {
		return errno.StatusParamsError, nil, err
	}

	if ret.Code != 0 {
		return int(ret.Code), nil, nil
	}

	data := map[string]interface{}{}
	data["token"] = ret.LoginResp.Token

	return errno.StatusOK, data, nil
}

func Register(username string, password string) (int, error) {
	conn, err := resource.CP.Get()
	defer resource.CP.Put(conn)

	if err != nil {
		log.Println(err.Error())
		return errno.StatusParamsError, err
	}

	session := rpc.NewSession(conn)

	msg := pb.Msg{ServiceName: "register", RegisterRequest: &pb.RegisterRequest{Username: username, Password: password}}

	ret, err := rpc.SendRequestForResp(session, &msg)

	if ret.Code != 0 {
		return int(ret.Code), err
	}

	return errno.StatusOK, nil
}

func ChangeNickname(username string, nickname string) (int, error) {
	conn, err := resource.CP.Get()
	defer resource.CP.Put(conn)

	if err != nil {
		log.Println(err.Error())
		return errno.StatusParamsError, err
	}

	session := rpc.NewSession(conn)

	msg := pb.Msg{ServiceName: "changenickname", ChangeNicknameRequest: &pb.ChangeNicknameRequest{Username: username, Nickname: nickname}}

	_, err = rpc.SendRequestForResp(session, &msg)

	if err != nil {
		return errno.StatusServerError, err
	}

	return errno.StatusOK, nil
}

func GetUserProfile(username string) (int, interface{}, error) {
	conn, err := resource.CP.Get()
	defer resource.CP.Put(conn)

	if err != nil {
		log.Println(err.Error())
		return errno.StatusServerError, nil, err
	}

	session := rpc.NewSession(conn)
	msg := pb.Msg{ServiceName: "getprofile", GetProfileRequest: &pb.GetProfileRequest{Username: username}}

	ret, err := rpc.SendRequestForResp(session, &msg)
	if err != nil {
		return errno.StatusServerError, nil, err
	}
	if ret.Code != 0 {
		return errno.StatusServerError, nil, nil
	}

	data := map[string]interface{}{
		"username": ret.GetProfileResp.Username,
		"nickname": ret.GetProfileResp.Nickname,
		"pic":      ret.GetProfileResp.Pic,
	}

	return errno.StatusOK, data, nil
}

func UpdatePic(username string, pic string) (int, error) {
	conn, err := resource.CP.Get()
	defer resource.CP.Put(conn)

	if err != nil {
		log.Println(err.Error())
		return errno.StatusParamsError, err
	}

	session := rpc.NewSession(conn)

	msg := pb.Msg{ServiceName: "updateprofile", UpdatePicRequest: &pb.UpdatePicRequest{Username: username, Pic: pic}}

	ret, err := rpc.SendRequestForResp(session, &msg)

	if ret.Code != 0 {
		return int(ret.Code), err
	}

	return errno.StatusOK, nil
}
