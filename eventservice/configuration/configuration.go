package configuration

import (
	"encoding/json"
	"log"
	"os"

	"github.com/agelloz/myevents/eventservice/persistence"
)

var (
	// DBTypeDefault is the default DB
	DBTypeDefault = persistence.DBType("mongodb")
	// DBConnectionDefault is the default connection
	DBConnectionDefault = "mongodb://127.0.0.1"
	// EndpointDefault is the default endpoint listening to HTTP
	EndpointDefault = "localhost:8181"
	// TLSEndpointDefault is the default endpoint listening to HTTPS
	TLSEndpointDefault = "localhost:9191"
	// AMPQURLDefault is the default url for the AMPQ broker
	AMPQURLDefault = "amqp://guest:guest@localhost:5672"
)

// ServiceConfig is
type ServiceConfig struct {
	DBType            persistence.DBType `json:"dbType"`
	DBConnection      string             `json:"dbConnection"`
	Endpoint          string             `json:"endpoint"`
	TLSEndpoint       string             `json:"tlsEndpoint"`
	AMQPMessageBroker string             `json:"amqp_message_broker"`
}

// ExtractConfiguration is
func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		EndpointDefault,
		TLSEndpointDefault,
		AMPQURLDefault,
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Println("configuration file not found. Continuing with default values:")
		log.Printf("%+v\n", conf)
		return conf, nil
	}
	err = json.NewDecoder(file).Decode(&conf)
	log.Printf("%+v\n", conf)
	return conf, err
}
