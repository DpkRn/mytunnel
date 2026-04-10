package config

import (
	"os"

	"github.com/joho/godotenv"
)

type TCPCfg struct {
	ListenAddr       string
	PublicURLScheme  string
	PublicHostSuffix string
}

func (c config) TCPServer() TCPCfg {
	_ = godotenv.Load()
	listenAddr := os.Getenv("TCP_LISTEN_ADDR")
	publicURLScheme := os.Getenv("PUBLIC_URL_SCHEME")
	publicHostSuffix := os.Getenv("PUBLIC_HOST_SUFFIX")

	switch c.TierFunc() {
	case "dev":
		return TCPCfg{
			ListenAddr:       listenAddr,
			PublicURLScheme:  "http",
			PublicHostSuffix: "localhost",
		}
	default:

		return TCPCfg{
			ListenAddr:       listenAddr,
			PublicURLScheme:  publicURLScheme,
			PublicHostSuffix: publicHostSuffix,
		}
	}
}

func (c TCPCfg) ListenAddrFunc() string {
	return c.ListenAddr
}

func (c TCPCfg) PublicURLSchemeFunc() string {
	return c.PublicURLScheme
}

func (c TCPCfg) PublicHostSuffixFunc() string {
	return c.PublicHostSuffix
}
