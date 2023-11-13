package utils

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	AppName string `yaml:"app_name"` //name of the app used as mqtt topic
}

type RemoteBrokerConfig struct {
	Broker      string `yaml:"broker"`       //url of the remote broker
	ClientId    string `yaml:"client_id"`    //client id to connect to remote broker
	Username    string `yaml:"username"`     //username to connect to remote broker
	Password    string `yaml:"password"`     //password to connect to remote broker
	BridgeTopic string `yaml:"bridge_topic"` //topic to communicate with the bridge
}

type MueckeConfig struct {
	RemoteBroker RemoteBrokerConfig `yaml:"remote_broker"`
	AppConfigs   []AppConfig        `yaml:"app_configs"`
}

func ParseConfig(file string) *MueckeConfig {
	// Read the yaml file
	cfg := MueckeConfig{}
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal the yaml file into the config struct
	err = yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Parse the flags
	flag.Parse()

	return &cfg
}
