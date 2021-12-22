package ws

import (
	"os"
	"strings"
)

// Config initialises and stores environment variables.
type Config struct {
	Stage            string
	ConnectionsTable string
	APIGatewayDomain string
}

func NewConfig() *Config {
	// ApiEndpoint attribute of AWS::ApiGatewayV2::Api includes protocol,
	// so this is stripped to get the domain.
	wsDomain := os.Getenv("API_GATEWAY_DOMAIN")
	wsDomain = strings.TrimPrefix(wsDomain, "wss://")

	return &Config{
		Stage:            os.Getenv("STAGE"),
		ConnectionsTable: os.Getenv("CONNECTIONS_TABLE"),
		APIGatewayDomain: wsDomain,
	}
}
