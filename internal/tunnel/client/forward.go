package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/DpkRn/devtunnel/internal/protocol"
	"github.com/hashicorp/yamux"
)

func acceptStreams(session *yamux.Session, localHost, port string) {
	for {
		stream, err := session.Accept()
		if err != nil {
			return
		}
		go handleStream(stream, localHost, port)
	}
}

func handleStream(stream net.Conn, localHost, port string) {
	defer stream.Close()

	reader := bufio.NewReader(stream)
	data, _ := reader.ReadBytes('\n')

	var req protocol.TunnelRequest
	json.Unmarshal(data, &req)

	httpReq, _ := http.NewRequest(
		req.Method,
		"http://"+localHost+":"+port+req.Path,
		bytes.NewReader(req.Body),
	)

	for k, v := range req.Headers {
		for _, val := range v {
			httpReq.Header.Add(k, val)
		}
	}

	resp, _ := http.DefaultClient.Do(httpReq)
	body, _ := io.ReadAll(resp.Body)

	response := protocol.TunnelResponse{
		Status:  resp.StatusCode,
		Headers: resp.Header,
		Body:    body,
	}

	out, _ := json.Marshal(response)
	stream.Write(append(out, '\n'))
}
