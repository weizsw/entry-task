package bootstrap

import (
	"database/sql"
	"log"
	"net"
	_ "net/http/pprof"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"

	"github.com/weizsw/entry-task/server/controller"
	"github.com/weizsw/entry-task/server/resource"
)

func Init() {
	initMysqlOnce()
	initRedisOnce()
	startServer()
}

func initMysqlOnce() {
	client, err := sql.Open("mysql", "root:CdB5f2vY@tcp(127.0.0.1:3306)/user?charset=utf8")
	if err != nil {
		panic(err)
	}

	client.SetMaxOpenConns(1000)
	client.SetMaxIdleConns(1000)
	resource.MysqlClient = client

	stmt, err := client.Prepare("SELECT * FROM users WHERE username = ?")
	if err != nil {
		panic(err)
	}
	resource.UserInfoStatement = stmt
}

func initRedisOnce() {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		PoolSize:     1000,
		MinIdleConns: 1000,
		MaxRetries:   2,
	})
	resource.RedisClient = rdb
}

func startServer() {
	listener, err := net.Listen("tcp", "127.0.0.1:8989")
	if err != nil {
		log.Fatal(err)
	}

	for {
		if conn, err := listener.Accept(); err == nil {
			go controller.HandleConn(conn)
		}
	}
}
