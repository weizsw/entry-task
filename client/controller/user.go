package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/weizsw/entry-task/client/errno"
	"github.com/weizsw/entry-task/client/models"
	"github.com/weizsw/entry-task/client/utils"
	"github.com/weizsw/entry-task/pb"
)

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}
	_ = t.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var reqBody *pb.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("login credentials error"))
		return
	}

	if len(reqBody.Username) == 0 {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("username is empty"))
		return
	}
	log.Println(reqBody.Username)
	if len(reqBody.Password) == 0 {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("password is empty"))
		return
	}

	code, data, err := models.Login(reqBody.Username, reqBody.Password)
	if err != nil {
		utils.RenderJson(w, code, nil, err)
		return
	}

	if code != errno.StatusOK {
		utils.RenderJson(w, code, nil, err)
		return
	}

	utils.RenderJson(w, errno.StatusOK, data, nil)
	return
}

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	if len(username) == 0 {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("username is empty"))
		return
	}

	password := r.PostFormValue("password")
	if len(password) == 0 {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("password is empty"))
		return
	}

	code, err := models.Register(username, password)

	utils.RenderJson(w, code, "ok", err)
}

func ChangeNickname(w http.ResponseWriter, r *http.Request) {
	reqBody := &pb.ChangeNicknameRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("change nickname error"))
		return
	}

	if len(reqBody.Username) == 0 || len(reqBody.Nickname) == 0 || len(reqBody.Token) == 0 {
		utils.RenderJson(w, errno.StatusParamsError, nil, nil)
		return
	}

	code, err := models.ChangeNickname(reqBody.Username, reqBody.Nickname, reqBody.Token)

	utils.RenderJson(w, code, nil, err)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	var reqBody *pb.GetProfileRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		utils.RenderJson(w, errno.StatusParamsError, nil, err)
		return
	}

	code, data, err := models.GetUserProfile(reqBody.Username, reqBody.Token)
	if err != nil {
		utils.RenderJson(w, code, nil, err)
		return
	}

	utils.RenderJson(w, code, data, err)
}

func UpdatePic(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	if len(username) == 0 {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("username is empty"))
		return
	}
	token := r.FormValue("token")
	if len(token) == 0 {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("token is empty"))
	}

	pic, header, err := r.FormFile("pic")
	if err != nil {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("please upload a pic"))
		return
	}

	ext := strings.ToLower(path.Ext(header.Filename))
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("jpg/png supported only"))
		return
	}

	url := "./assets/" + username + ext
	saveFile, err := os.OpenFile(url, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(saveFile, pic)

	defer pic.Close()
	defer saveFile.Close()

	url = "http://localhost:8080/assets/" + username + ext
	_, err = models.UpdatePic(username, url, token)

	ret := map[string]interface{}{}
	ret["pic"] = url

	utils.RenderJson(w, errno.StatusOK, ret, err)
}

func AssetsFunc(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./assets"))
	http.StripPrefix("/assets/", fs)
	// http.Handle("/assets/", fs)
	// http.ServeFile(w, r, "assets/default.jpeg")
}
