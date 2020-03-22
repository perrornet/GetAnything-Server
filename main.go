package main

import (
	"GetAnything-Server/server"
	"GetAnything-Server/utils"
)

func main() {
	// 开发时需要注释
	//u, err := update.NewUpdate("https://GetAnything-Server")
	//if err != nil{
	//	log.Panic(err)
	//}
	//go u.Restart()
	// end
	server.StartServer(utils.GetCmd())
}
