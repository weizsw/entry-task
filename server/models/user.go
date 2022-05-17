package models

import (
	"fmt"
	"log"

	"github.com/weizsw/entry-task/client/errno"
	"github.com/weizsw/entry-task/pb"
	"github.com/weizsw/entry-task/rpc"
	db "github.com/weizsw/entry-task/server/dao/db"
	"github.com/weizsw/entry-task/server/dao/redis"
	"github.com/weizsw/entry-task/server/resource"
	"github.com/weizsw/entry-task/server/utils"
)

func Login(session *rpc.Session, req *pb.Msg) {
	msg := &pb.Msg{}

	res, err := db.GetProfileWithPrepare(resource.UserProfileStatement, req.LoginRequest.Username)
	if err != nil {
		log.Println(err.Error())
		msg = &pb.Msg{Code: errno.LoginCredentialError, Message: err.Error()}
		rpc.SendRequest(session, msg)
		return
	}

	if res.Password != req.LoginRequest.Password {
		msg = &pb.Msg{Code: errno.LoginCredentialError, Message: errno.ErrorMsg[errno.LoginCredentialError]}
		rpc.SendRequest(session, msg)
		return
	}

	token := utils.GetMD5Hash(req.LoginRequest.Username + req.LoginRequest.Password)
	err = redis.HSet(resource.RedisClient, fmt.Sprintf(redis.USER_TOKEN_FMT, req.LoginRequest.Username), token, 1)
	if err != nil {
		msg = &pb.Msg{Code: errno.LoginCredentialError, Message: errno.ErrorMsg[errno.LoginCredentialError]}
		rpc.SendRequest(session, msg)
		return
	}

	msg = &pb.Msg{Code: errno.StatusOK, Token: token, Message: errno.ErrorMsg[errno.StatusOK]}
	rpc.SendRequest(session, msg)
	return
}

func Register(session *rpc.Session, ret *pb.Msg) {
	hash := utils.GetMD5Hash(ret.GetRegisterRequest().Password)

	err := db.InsertUser(resource.MysqlClient, ret.GetRegisterRequest().Username, string(hash))
	if err != nil {
		log.Println(err)
		rpc.SendRequest(session, &pb.Msg{Response: &pb.Response{Code: errno.RegisterFailedError, Msg: err.Error()}})
	}
	msg := &pb.Msg{Response: &pb.Response{Code: errno.StatusOK, Msg: errno.ErrorMsg[errno.StatusOK]}}
	rpc.SendRequest(session, msg)
}

func GetUserProfile(session *rpc.Session, req *pb.Msg) {
	msg := &pb.Msg{}

	_, err := redis.HGet(resource.RedisClient, fmt.Sprintf(redis.USER_TOKEN_FMT, req.GetProfileRequest.Username), req.Token)
	if err != nil {
		log.Println(err.Error())
		msg = &pb.Msg{Code: errno.MySQLError, Message: err.Error()}
		rpc.SendRequest(session, msg)
		return
	}

	cache, err := redis.Get(resource.RedisClient, req.GetProfileRequest.Username)
	if err != nil {
		log.Println(err.Error())
		res, err := db.GetUser(resource.MysqlClient, req.GetProfileRequest.Username)
		if err != nil {
			msg = &pb.Msg{Code: errno.MySQLError, Message: err.Error()}
			rpc.SendRequest(session, msg)
			return
		}

		msg = &pb.Msg{Code: errno.StatusOK, GetProfileResp: &pb.GetProfileResp{Username: res.Username, Nickname: res.Nickname, Pic: res.ProfilePic, Password: res.Password}, Message: errno.ErrorMsg[errno.StatusOK]}

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

func ChangeNickname(session *rpc.Session, req *pb.Msg) {
	msg := &pb.Msg{}
	_, err := redis.HGet(resource.RedisClient, fmt.Sprintf(redis.USER_TOKEN_FMT, req.ChangeNicknameRequest.Username), req.Token)
	if err != nil {
		log.Println(err.Error())
		msg = &pb.Msg{Code: errno.RedisError, Message: err.Error()}
		rpc.SendRequest(session, msg)
		return
	}

	err = db.UpdateNickname(resource.MysqlClient, req.ChangeNicknameRequest.Username, req.ChangeNicknameRequest.Nickname)
	if err != nil {
		log.Println(err.Error())
		msg = &pb.Msg{Code: errno.MySQLError, Message: err.Error()}
		rpc.SendRequest(session, msg)
		return
	}

	msg = &pb.Msg{Code: errno.StatusOK, Message: errno.ErrorMsg[errno.StatusOK]}
	err = redis.Del(resource.RedisClient, req.ChangeNicknameRequest.Username)
	if err != nil {
		msg = &pb.Msg{Code: errno.RedisError, Message: err.Error()}
	}

	rpc.SendRequest(session, msg)
	return
}

func UpdatePic(session *rpc.Session, req *pb.Msg) {
	msg := &pb.Msg{}
	_, err := redis.HGet(resource.RedisClient, fmt.Sprintf(redis.USER_TOKEN_FMT, req.UpdateProfileRequest.Username), req.Token)
	if err != nil {
		log.Println(err.Error())
		msg = &pb.Msg{Code: errno.RedisError, Message: err.Error()}
		rpc.SendRequest(session, msg)
		return
	}

	err = db.UpdateProfile(resource.MysqlClient, req.UpdateProfileRequest.Username, req.UpdateProfileRequest.Pic)
	if err != nil {
		log.Fatal(err.Error())
	}

	msg = &pb.Msg{Response: &pb.Response{Code: errno.StatusOK, Msg: errno.ErrorMsg[errno.StatusOK]}}

	err = redis.Del(resource.RedisClient, req.UpdateProfileRequest.Username)
	if err != nil {
		msg = &pb.Msg{Response: &pb.Response{Code: errno.StatusServerError, Msg: errno.ErrorMsg[errno.StatusServerError]}}
	}
	rpc.SendRequest(session, msg)
}
