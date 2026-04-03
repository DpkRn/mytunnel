package config

// Defaults for the tunnel stack. Centralize here so CLI, server, and library stay in sync.

const (
	DefaultTunnelServerAddr = "localhost:9000"

	ControlListenAddr = ":9000"
	HTTPListenAddr    = ":3000"

	// PublicHostSuffix is appended after the subdomain (must match edge HTTP server).
	PublicHostSuffix = ".localhost:3000"
)
