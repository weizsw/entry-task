package models

import (
	"log"

	"github.com/weizsw/entry-task/client/errno"
	"github.com/weizsw/entry-task/pb"
	"github.com/weizsw/entry-task/rpc"
	db "github.com/weizsw/entry-task/server/dao/db"
	"github.com/weizsw/entry-task/server/dao/redis"
	"github.com/weizsw/entry-task/server/resource"
	"github.com/weizsw/entry-task/server/utils"
)

func Login(session *rpc.Session, ret *pb.Msg) {
	var msg *pb.Msg
	msg = &pb.Msg{Response: &pb.Response{Code: errno.StatusOK, Msg: errno.ErrorMsg[errno.StatusOK]}}

	res, err := db.GetUserWithPrepare(resource.UserInfoStatement, ret.GetLoginRequest().Username)
	// res, err := db.GetUser(resource.MysqlClient, ret.GetLoginRequest().Username)
	if err != nil {
		log.Println(err.Error())
		msg = &pb.Msg{Response: &pb.Response{Code: errno.LoginCredentialError, Msg: err.Error()}}
		rpc.SendRequest(session, msg)
		return
	}

	hash := utils.GetMD5Hash(ret.GetLoginRequest().Password)
	if res.Password != hash {
		msg = &pb.Msg{Response: &pb.Response{Code: errno.LoginCredentialError, Msg: "password error"}}
		rpc.SendRequest(session, msg)
		return
	}

	// err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(ret.GetLoginRequest().Password))
	// if err != nil {
	// 	msg = &pb.Msg{Response: &pb.Response{Code: errno.LoginCredentialError, Msg: err.Error()}}
	// }

	rpc.SendRequest(session, msg)
	return
}

func Register(session *rpc.Session, ret *pb.Msg) {

	hash := utils.GetMD5Hash(ret.GetRegisterRequest().Password)
	// hash, err := bcrypt.GenerateFromPassword([]byte(ret.GetRegisterRequest().Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	err := db.InsertUser(resource.MysqlClient, ret.GetRegisterRequest().Username, string(hash))
	if err != nil {
		log.Println(err)
		rpc.SendRequest(session, &pb.Msg{Response: &pb.Response{Code: errno.RegisterFailedError, Msg: err.Error()}})
	}
	msg := &pb.Msg{Response: &pb.Response{Code: errno.StatusOK, Msg: errno.ErrorMsg[errno.StatusOK]}}
	rpc.SendRequest(session, msg)
}

func UserInfo(session *rpc.Session, ret *pb.Msg) {
	msg := &pb.Msg{}
	cache, err := redis.Get(resource.RedisClient, ret.GetUserInfoRequest().Username)
	if err != nil {
		res, err := db.GetUser(resource.MysqlClient, ret.UserInfoRequest.GetUsername())
		if err != nil {
			msg = &pb.Msg{Response: &pb.Response{Code: errno.LoginCredentialError, Msg: err.Error()}}
			rpc.SendRequest(session, msg)
			return
		}

		msg = &pb.Msg{UserInfoResponse: &pb.UserInfoResponse{Username: res.Username, Nickname: res.Nickname, Pic: res.ProfilePic, Password: res.Password}}

		go func() {
			toCache, err := msg.Encode()
			if err != nil {
				log.Println(err)
			}
			redis.Set(resource.RedisClient, res.Username, toCache)
		}()
	} else {
		err = msg.Decode([]byte(cache))
		if err != nil {
			log.Println(err)
		}
	}

	rpc.SendRequest(session, msg)
	return
}

func ChangeNickname(session *rpc.Session, ret *pb.Msg) {
	err := db.UpdateNickname(resource.MysqlClient, ret.ChangeNicknameRequest.GetUsername(), ret.ChangeNicknameRequest.GetNickname())
	if err != nil {
		log.Fatal(err.Error())
	}
	msg := &pb.Msg{Response: &pb.Response{Code: errno.StatusOK, Msg: errno.ErrorMsg[errno.StatusOK]}}

	err = redis.Del(resource.RedisClient, ret.ChangeNicknameRequest.GetUsername())
	if err != nil {
		msg = &pb.Msg{Response: &pb.Response{Code: errno.StatusServerError, Msg: errno.ErrorMsg[errno.StatusServerError]}}
	}
	rpc.SendRequest(session, msg)
	return
}

func UpdateProfile(session *rpc.Session, ret *pb.Msg) {
	err := db.UpdateProfile(resource.MysqlClient, ret.UpdateProfileRequest.GetUsername(), ret.UpdateProfileRequest.GetPic())
	if err != nil {
		log.Fatal(err.Error())
	}
	msg := &pb.Msg{Response: &pb.Response{Code: errno.StatusOK, Msg: errno.ErrorMsg[errno.StatusOK]}}

	err = redis.Del(resource.RedisClient, ret.UpdateProfileRequest.GetUsername())
	if err != nil {
		msg = &pb.Msg{Response: &pb.Response{Code: errno.StatusServerError, Msg: errno.ErrorMsg[errno.StatusServerError]}}
	}
	rpc.SendRequest(session, msg)
}
