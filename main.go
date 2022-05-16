package main

import (
	"log"
	"os"
	"runtime/trace"
	"sync"
	"time"

	client "github.com/weizsw/entry-task/client/bootstrap"
	server "github.com/weizsw/entry-task/server/bootstrap"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		log.Println("initializing server...")
		server.Init()
		wg.Done()
	}()
	time.Sleep(time.Second * 2)
	go func() {
		log.Println("initializing client...")
		client.Init()
		wg.Done()
	}()

	wg.Wait()

	trace.Start(os.Stderr)
	defer trace.Stop()

	ch := make(chan string)
	go func() {
		ch <- "EDDYCJY"
	}()

	<-ch
}
