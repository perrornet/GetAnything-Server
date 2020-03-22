package utils

import (
	"os"
	"strings"
)

type CmdArgs struct {
	Host string
	Port string
	Mode string
}

func GetCmd() *CmdArgs {
	cmd := &CmdArgs{}
	cmd.Host = os.Getenv("HOST")
	if strings.TrimSpace(cmd.Host) == "" {
		cmd.Host = "0.0.0.0"
	}
	cmd.Port = os.Getenv("PORT")
	if strings.TrimSpace(cmd.Port) == "" {
		cmd.Port = "80"
	}
	mode := os.Getenv("MODE")
	switch strings.TrimSpace(mode) {
	case "debug", "release":
		cmd.Mode = strings.TrimSpace(mode)
	default:
		cmd.Mode = "release"

	}
	return cmd
}
