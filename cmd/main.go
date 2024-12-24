package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	_ "net/http/pprof"

	"github.com/xhuliodo/p2p_video_chat_backend/config"
	"github.com/xhuliodo/p2p_video_chat_backend/server"
)

func main() {
	conf := config.NewConfig([]string{})
	server := server.NewServer(conf)

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	// ticker := time.NewTicker(5 * time.Second)
	// quit := make(chan struct{})
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			var m runtime.MemStats
	// 			runtime.ReadMemStats(&m)
	// 			metrics := fmt.Sprintf(
	// 				"Alloc: %.2f MB, TotalAlloc: %.2f MB, Sys: %.2f MB, NumGC: %d\n",
	// 				float64(m.Alloc)/1024/1024,
	// 				float64(m.TotalAlloc)/1024/1024,
	// 				float64(m.Sys)/1024/1024,
	// 				m.NumGC,
	// 			)

	// 			log.Println(metrics)
	// 		case <-quit:
	// 			ticker.Stop()
	// 			log.Println("stopping logging memory usage")
	// 			return
	// 		}
	// 	}
	// }()

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := server.Start(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// quit <- struct{}{}

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), conf.Server.GracefulTimeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
