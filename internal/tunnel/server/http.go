package server

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/DpkRn/devtunnel/internal/protocol"
)

// EdgeHTTP serves public HTTP and routes by subdomain to the correct tunnel session.
type EdgeHTTP struct {
	Registry *Registry
}

func NewEdgeHTTP(reg *Registry) *EdgeHTTP {
	return &EdgeHTTP{Registry: reg}
}

// Handler returns the root http.HandlerFunc for the edge server.
func (e *EdgeHTTP) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.Host, ".")
		if len(parts) < 2 {
			http.Error(w, "Invalid host", http.StatusBadRequest)
			return
		}
		sub := parts[0]

		session, ok := e.Registry.Get(sub)
		if !ok {
			http.Error(w, "Tunnel not found", http.StatusNotFound)
			return
		}

		stream, err := session.Open()
		if err != nil {
			e.Registry.Remove(sub)
			http.Error(w, "Failed to open stream", http.StatusInternalServerError)
			return
		}
		defer stream.Close()

		body, _ := io.ReadAll(r.Body)

		req := protocol.TunnelRequest{
			Method:  r.Method,
			Path:    r.URL.String(),
			Headers: r.Header,
			Body:    body,
		}

		data, err := json.Marshal(req)
		if err != nil {
			e.Registry.Remove(sub)
			http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
			return
		}
		if _, err := stream.Write(append(data, '\n')); err != nil {
			e.Registry.Remove(sub)
			http.Error(w, "Failed to write request", http.StatusInternalServerError)
			return
		}

		reader := bufio.NewReader(stream)
		respBytes, err := reader.ReadBytes('\n')
		if err != nil {
			e.Registry.Remove(sub)
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		var resp protocol.TunnelResponse
		if err := json.Unmarshal(respBytes, &resp); err != nil {
			http.Error(w, "Invalid tunnel response", http.StatusBadGateway)
			return
		}

		for k, v := range resp.Headers {
			for _, val := range v {
				w.Header().Add(k, val)
			}
		}
		w.WriteHeader(resp.Status)
		_, _ = w.Write(resp.Body)
	}
}
