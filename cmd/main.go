package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/xhuliodo/p2p_video_chat_backend/config"
	"github.com/xhuliodo/p2p_video_chat_backend/server"
)

func main() {
	conf := config.NewConfig()
	server := server.NewServer(conf)

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
