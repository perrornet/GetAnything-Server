package utils

import "flag"

type CmdArgs struct {
	Host *string
	Port *string
	Mode *string
}

func GetCmd() *CmdArgs {
	cmd := &CmdArgs{}
	cmd.Host = flag.String("h", "0.0.0.0", "限制外部访问IP地址（默认：0.0.0.0）")
	cmd.Port = flag.String("p", "80", "监听端口（默认：80）")
	cmd.Mode = flag.String("m", "release", "当前服务器状态（默认： release， 可选：debug/release）")
	flag.Parse()
	return cmd
}
