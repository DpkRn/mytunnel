package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/DpkRn/devtunnel/internal/protocol"
)

func Handler(reg *Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("r.Host:", r.Host)
		parts := strings.Split(r.Host, ".")
		if len(parts) < 2 {
			http.Error(w, "Invalid host", 400)
			return
		}

		subdomain := parts[0]

		session, ok := reg.Get(subdomain)
		if !ok {
			http.Error(w, "Tunnel not found", 404)
			return
		}

		stream, _ := session.Open()
		defer stream.Close()

		body, _ := io.ReadAll(r.Body)

		req := protocol.TunnelRequest{
			Method:  r.Method,
			Path:    r.URL.String(),
			Headers: r.Header,
			Body:    body,
		}

		fmt.Println("req:", req)

		data, _ := json.Marshal(req)
		stream.Write(append(data, '\n'))

		reader := bufio.NewReader(stream)
		respBytes, _ := reader.ReadBytes('\n')
		fmt.Println("respBytes:", string(respBytes))

		var resp protocol.TunnelResponse
		json.Unmarshal(respBytes, &resp)

		for k, v := range resp.Headers {
			for _, val := range v {
				w.Header().Add(k, val)
			}
		}

		w.WriteHeader(resp.Status)
		w.Write(resp.Body)
	}
}
