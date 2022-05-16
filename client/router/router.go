package router

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/weizsw/entry-task/client/controller"
)

func Init() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:7000", nil))
	// }()
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", controller.Index)
	mux.HandleFunc("/register", controller.Register)
	mux.HandleFunc("/login", controller.Login)
	mux.HandleFunc("/nickname", controller.ChangeNickname)
	mux.HandleFunc("/pic", controller.UpdateProfile)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))
}
