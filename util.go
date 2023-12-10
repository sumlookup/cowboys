package main

import "os"

const (
	DEFAULT_PORT = "9090"
	SERVICE_NAME = "cowboys"
	TRANSPORT    = "grpc"
	REGISTRY     = "mdns"
	SELECTOR     = "registry"

	MODE = "intermediate"
)

func GetPort() string {
	port := DEFAULT_PORT
	value := os.Getenv("HTTP_PORT")
	if len(value) == 0 {
		return port
	}
	return value
}

func GetSelector() string {
	selector := SELECTOR
	value := os.Getenv("SELECTOR")
	if len(value) == 0 {
		return selector
	}
	return value
}

func GetTransport() string {
	transport := TRANSPORT
	value := os.Getenv("TRANSPORT")
	if len(value) == 0 {
		return transport
	}
	return value
}

func GetRegistry() string {
	registry := REGISTRY
	value := os.Getenv("REGISTRY")
	if len(value) == 0 {
		return registry
	}
	return value
}

func GetGameMode() string {
	mode := MODE
	value := os.Getenv("MODE")
	if len(value) == 0 {
		return mode
	}
	return value
}

func GetServiceName() string {
	svc := SERVICE_NAME
	value := os.Getenv("SERVICE_NAME")
	if len(value) == 0 {
		return svc
	}
	return value
}
