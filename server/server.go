package server

import (
	"GetAnything-Server/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartServer(args *utils.CmdArgs) {
	gin.SetMode(args.Mode)
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", args.Host, args.Port),
		Handler: routers(),
	}
	if err := Cache.LoadFile("./GetAnything.cache"); err != nil {
		log.Println(err)
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	log.Println("Welcome to GetAnything service. Your server address is:", fmt.Sprintf("%s:%s", args.Host, args.Port))
	<-quit
	if err := Cache.SaveFile("./GetAnything.cache"); err != nil {
		log.Println(err)
	}
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
