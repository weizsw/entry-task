package router

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/weizsw/entry-task/client/controller"
)

func Init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:7000", nil))
	}()
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/register", controller.Register)
	mux.HandleFunc("/login", Middleware(http.HandlerFunc(controller.Login)))
	mux.HandleFunc("/nickname", Middleware(http.HandlerFunc(controller.ChangeNickname)))
	mux.HandleFunc("/pic", Middleware(http.HandlerFunc(controller.UpdatePic)))
	mux.HandleFunc("/profile", Middleware(http.HandlerFunc(controller.GetProfile)))

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))
}

func Middleware(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
	})
}
