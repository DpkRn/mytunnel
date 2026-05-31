package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/DpkRn/devtunnel/internal/pkg"
	"github.com/DpkRn/devtunnel/internal/platform/config"
	"github.com/DpkRn/devtunnel/internal/platform/mongo"
	"github.com/hashicorp/yamux"
)

type TunnelType string

const (
	TunnelTypeGotunnel   TunnelType = "gotunnel"
	TunnelTypeNodeTunnel TunnelType = "nodetunnel"
)

type ConnectionStatus string

const (
	ConnectionStatusOnline  ConnectionStatus = "online"
	ConnectionStatusOffline ConnectionStatus = "offline"
)

type TunnelRequest struct {
	TunnelType       TunnelType       `json:"tunnel_type"`
	Version          string           `json:"version"`
	ClientIP         string           `json:"client_ip"`
	ConnectionTime   time.Time        `json:"connection_time"`
	ConnectionType   string           `json:"connection_type"`
	ConnectionID     string           `json:"connection_id"`
	ConnectionStatus ConnectionStatus `json:"connection_status"`
}

func StartTCP(reg *Registry, tcpConfig config.TCPCfg, mongoClient mongo.Client) {
	// config := config.NewConfig()
	tcpPort := tcpConfig.ListenAddrFunc()
	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", tcpPort, err)
	}
	fmt.Println("✅TCP Connection Listening on port: ", tcpPort)
	for {
		conn, _ := listener.Accept()
		fmt.Println("Client connected:", conn.RemoteAddr())

		//request client welcome message
		reader := bufio.NewReader(conn)
		message, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatalf("Failed to read welcome message: %v", err)
		}
		tunnelReq := TunnelRequest{}
		err = json.Unmarshal(message, &tunnelReq)
		if err != nil {
			log.Fatalf("Failed to marshal tunnel request: %v", err)
		}
		fmt.Println("tunnelReq:", tunnelReq)
		tunnelReq.ClientIP = conn.RemoteAddr().String()
		tunnelReq.ConnectionTime = time.Now()
		tunnelReq.ConnectionType = "tcp"
		tunnelReq.ConnectionStatus = ConnectionStatusOnline
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("panic:", r)
				}
			}()
			mongoClient.InsertTunnelLog(context.Background(), tunnelReq)
		}()
		go func() {
			handleClient(conn, reg, tcpConfig)
		}()
	}
}

func handleClient(conn net.Conn, reg *Registry, tcpConfig config.TCPCfg) {
	subdomain := pkg.GenerateID()
	scheme := strings.TrimSuffix(strings.ToLower(tcpConfig.PublicURLSchemeFunc()), "://")
	if scheme == "" {
		scheme = "https"
	}
	publicURL := fmt.Sprintf("%s://%s.%s\n", scheme, subdomain, tcpConfig.PublicHostSuffixFunc())

	// ✅ send BEFORE yamux
	conn.Write([]byte(publicURL))

	// now start yamux
	session, err := yamux.Server(conn, nil)
	if err != nil {
		conn.Close()
		return
	}

	reg.Add(subdomain, session)

	fmt.Println("Client connected:", subdomain)
}
