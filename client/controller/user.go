package controller

import (
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
	username := r.FormValue("username")
	log.Println(username)
	if len(username) == 0 {
		// utils.RenderIndex(w, errno.StatusParamsError)
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("username is empty"))
		return
	}

	password := r.FormValue("password")
	if len(password) == 0 {
		// utils.RenderIndex(w, errno.StatusParamsError)
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("password is empty"))
		return
	}

	// cookie, _ := r.Cookie(username)
	// if cookie != nil && cookie.Value == username {
	// 	// utils.RenderLogin(w, username, nil)
	// 	utils.RenderJson(w, 0, "ok", nil)
	// 	return
	// }

	code, err := models.Login(username, password)
	if err != nil {
		utils.RenderIndex(w, code)
		return
	}

	if code != errno.StatusOK {
		utils.RenderIndex(w, code)
		// utils.SetCookie(w, username, username)
	}

	utils.RenderLogin(w, username, nil)
	// utils.RenderJson(w, code, "ok", err)
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
	username := r.PostFormValue("username")
	if len(username) == 0 {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("username is empty"))
		return
	}

	nickname := r.PostFormValue("nickname")
	if len(nickname) == 0 {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("nickname is empty"))
		return
	}

	_, err := models.ChangeNickname(username, nickname)

	utils.RenderLogin(w, username, err)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	if len(username) == 0 {
		utils.RenderJson(w, errno.StatusParamsError, nil, errors.New("username is empty"))
		return
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

	path := "./assets/" + username + ext
	saveFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(saveFile, pic)

	defer pic.Close()
	defer saveFile.Close()

	_, err = models.UpdateProfile(username, path)
	utils.RenderLogin(w, username, err)
}

func AssetsFunc(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./assets"))
	http.StripPrefix("/assets/", fs)
	// http.Handle("/assets/", fs)
	// http.ServeFile(w, r, "assets/default.jpeg")
}
