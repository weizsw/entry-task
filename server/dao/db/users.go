package db

import (
	"database/sql"
	"log"
)

const (
	DEFAULT_PROFILE = "./assets/default.jpeg"
)

type User struct {
	ID         uint64 `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Password   string `json:"password"`
	ProfilePic string `json:"profile_pic"`
}

func GetUser(db *sql.DB, username string) (*User, error) {
	ret := &User{}
	err := db.QueryRow("Select * from users where username = ?", username).Scan(&ret.ID, &ret.Username, &ret.Nickname, &ret.Password, &ret.ProfilePic)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func GetUserWithPrepare(stmt *sql.Stmt, username string) (*User, error) {
	ret := &User{}

	err := stmt.QueryRow(username).Scan(&ret.ID, &ret.Username, &ret.Nickname, &ret.Password, &ret.ProfilePic)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return ret, nil
}

func InsertUser(db *sql.DB, username string, password string) error {
	_, err := db.Exec("insert into users (username, nickname, password, profile_pic) values(?,?,?,?)", username, "", password, DEFAULT_PROFILE)
	if err != nil {
		return err
	}

	return nil
}

func UpdateNickname(db *sql.DB, username string, nickname string) error {
	_, err := db.Exec("update users set nickname = ? where username = ?", nickname, username)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProfile(db *sql.DB, username string, pic string) error {
	_, err := db.Exec("update users set profile_pic = ? where username = ?", pic, username)
	if err != nil {
		return err
	}

	return nil
}
