package server

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/DpkRn/devtunnel/internal/pkg"
	"github.com/DpkRn/devtunnel/internal/platform/mongo"
	"github.com/DpkRn/devtunnel/internal/protocol"
)

// HandleTunnelRequest proxies an HTTP request through the yamux tunnel for tunnelID.
func HandleTunnelRequest(
	w http.ResponseWriter,
	r *http.Request,
	reg *Registry,
	mongoClient mongo.Client,
	tunnelID string,
) {
	start := time.Now()
	reqID := pkg.GenerateRequestID()

	fmt.Println("r.Host:", r.Host, "tunnelID:", tunnelID)

	session, ok := reg.Get(tunnelID)
	if !ok {
		http.Error(w, "Tunnel not found", http.StatusNotFound)
		return
	}

	stream, err := session.OpenStream()
	if err != nil {
		reg.Remove(tunnelID)
		http.Error(w, "Tunnel session closed", http.StatusBadGateway)
		return
	}
	streamID := stream.StreamID()
	defer stream.Close()

	body, _ := io.ReadAll(r.Body)

	req := protocol.TunnelRequest{
		Method:  r.Method,
		Path:    r.URL.String(),
		Headers: r.Header,
		Body:    body,
	}

	fmt.Println("req:", req)

	data, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusInternalServerError)
		return
	}
	if _, err := stream.Write(append(data, '\n')); err != nil {
		reg.Remove(tunnelID)
		http.Error(w, "Tunnel write failed", http.StatusBadGateway)
		return
	}

	reader := bufio.NewReader(stream)
	respBytes, err := reader.ReadBytes('\n')
	if err != nil || len(respBytes) == 0 {
		reg.Remove(tunnelID)
		http.Error(w, "Tunnel closed before response", http.StatusBadGateway)
		return
	}
	fmt.Println("respBytes:", string(respBytes))

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

	go func() {

		defer func() {
			if r := recover(); r != nil {
				log.Println("panic:", r)
			}
		}()

		doc := map[string]any{
			"stream_id":  streamID,
			"request_id": reqID,
			"tunnel_id":  tunnelID,

			"request": map[string]any{
				"method":  r.Method,
				"path":    r.URL.String(),
				"headers": r.Header,
				"body":    base64.StdEncoding.EncodeToString(body),
				"size":    len(body),
			},

			"response": map[string]any{
				"status":  resp.Status,
				"headers": resp.Headers,
				"body":    base64.StdEncoding.EncodeToString(resp.Body),
				"size":    len(resp.Body),
			},

			"timing": map[string]any{
				"duration_ms": time.Since(start).Milliseconds(),
				"timestamp":   start.Format(time.RFC3339),
			},
			"host": map[string]any{
				"client_ip":  r.RemoteAddr,
				"referrer":   r.Referer(),
				"user_agent": r.UserAgent(),
				"host":       r.Host,
			},

			"created_at": time.Now(),
		}

		_, err := mongoClient.InsertRequestLog(context.Background(), doc)
		if err != nil {
			log.Println("Mongo error:", err)
		}
	}()
}

func serveControlPlane(w http.ResponseWriter, r *http.Request, reg *Registry, mongoClient mongo.Client) {
	path := r.URL.Path
	switch {
	case path == "/health":
		HealthHandler(w, r)
	case path == "/logs":
		GetLogsHandler(mongoClient).ServeHTTP(w, r)
	case strings.HasPrefix(path, "/logs/"):
		GetLogByIDHandler(mongoClient).ServeHTTP(w, r)
	case strings.HasPrefix(path, "/replay/"):
		ReplayHandler(reg, mongoClient).ServeHTTP(w, r)
	case path == "/":
		ServerHomeHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}
