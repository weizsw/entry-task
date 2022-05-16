package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/weizsw/entry-task/client/errno"
	"github.com/weizsw/entry-task/client/models"
	"github.com/weizsw/entry-task/pb"
)

func NewSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return string(hash[:])
}

type JsonResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type UserInfoResponse struct {
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	ProfilePic string `json:"profile_pic"`
}

func RenderJson(w http.ResponseWriter, code int, data interface{}, err error) {
	w.Header().Set("content-type", "text/json")
	var msg string
	if err != nil {
		msg = err.Error()
	} else {
		msg = errno.ErrorMsg[code]
	}

	jsonRet := JsonResponse{Code: code, Data: data, Msg: msg}
	retMsg, _ := json.Marshal(jsonRet)
	if code < 300 {
		w.WriteHeader(200)
	} else if code < 500 {
		w.WriteHeader(400)
	} else {
		w.WriteHeader(500)
	}

	w.Write(retMsg)
}

func RenderLogin(w http.ResponseWriter, username string, err error) {
	_, data, err := models.GetUserInfo(username)
	if err != nil {
		log.Println(err)
		RenderJson(w, errno.StatusServerError, nil, err)
		return
	}
	t, err := template.ParseFiles("login.html")
	if err != nil {
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	_ = t.Execute(w, data.UserInfoResponse)
}

func RenderIndex(w http.ResponseWriter, code int) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}
	data := &pb.Msg{}
	data.Response = &pb.Response{Code: int64(code), Msg: errno.ErrorMsg[int(code)]}
	_ = t.Execute(w, data.Response)
}

func SetCookie(w http.ResponseWriter, name string, value string) {
	expiration := time.Now()
	expiration = expiration.AddDate(0, 0, 1)
	cookie := http.Cookie{Name: name, Value: value, Expires: expiration}
	http.SetCookie(w, &cookie)
}
