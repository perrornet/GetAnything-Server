package main

import (
	"github.com/PerrorOne/GetAnything-Server/server"
)

func main() {
	// 开发时需要注释
	//u, err := update.NewUpdate("https://github.com/PerrorOne/GetAnything-Server")
	//if err != nil{
	//	log.Panic(err)
	//}
	//go u.Restart()
	// end
	server.StartServer("", "8080")
}
