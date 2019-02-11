package main

import "github.com/PerrorOne/GetAnything-Server/server"

func main() {
	// 开发时需要注释
	//u, err := update.NewUpdate("https://github.com/PerrorOne/GetAnything-Server")
	//if err != nil{
	//	log.Panic(err)
	//}
	//go u.Restart()
	// end
	server.StartServer("", "8080")

	//fmt.Println(url.QueryEscape("https://golangtc.com/t/55d6e5f7b09ecc2e87000009/尼玛死了"))
	//d, _ := extractors.Match("https://v.douyu.com/show/Qyz171QQB1G7BJj9")
	//fmt.Println(d.GetFileInfo())
}
