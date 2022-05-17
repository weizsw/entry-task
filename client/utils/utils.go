package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/weizsw/entry-task/client/errno"
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

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func RenderJson(w http.ResponseWriter, code int, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
		if code == 401 {
			w.WriteHeader(401)
		} else {
			w.WriteHeader(400)
		}
	} else {
		w.WriteHeader(500)
	}

	w.Write(retMsg)
}

func SetCookie(w http.ResponseWriter, name string, value string) {
	expiration := time.Now()
	expiration = expiration.AddDate(0, 0, 1)
	cookie := http.Cookie{Name: name, Value: value, Expires: expiration}
	http.SetCookie(w, &cookie)
}
