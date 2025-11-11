package main

// Version information
// These variables can be overridden at build time using ldflags:
// go build -ldflags "-X main.Version=v1.0.0 -X main.GitCommit=abc123"
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)
