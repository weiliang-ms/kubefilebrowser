package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"kubecp/configs"
	_ "kubecp/configs"
	"kubecp/logs"
	"kubecp/routers"
	_ "kubecp/routers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title KubeCp Swagger
// @version 1.0
// @description 网页版kubectl cp
// @BasePath /
// @query.collection.format multi
func main() {
	gin.SetMode(configs.Config.RunMode)
	logs.Info("Start up...")
	r := routers.Router()
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", configs.Config.HTTPAddr, configs.Config.HTTPPort),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 300,
		ReadTimeout:  time.Second * 300,
		IdleTimeout:  time.Second * 300,
		Handler:      r, // Pass our instance of gin in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logs.Error(err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logs.Error("shutting down")
	os.Exit(0)
}
