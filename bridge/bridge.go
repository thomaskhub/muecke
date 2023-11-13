package bridge

import (
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thomaskhub/muecke/client"
	"github.com/thomaskhub/muecke/server"
	"github.com/thomaskhub/muecke/utils"
)

//parse flag to the config file

func StartBridge(cfg *utils.MueckeConfig) {

	client := client.NewMqttClientWithConfig(
		cfg.RemoteBroker.Broker,
		cfg.RemoteBroker.ClientId,
		cfg.RemoteBroker.Username,
		cfg.RemoteBroker.Password,
	)

	server := server.MochiServer{}
	server.StartServer(cfg, client)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// now subscibe with the paho client to the bridge topic, can only be done
	// after the server is running
	remoteToLocal := func(client mqtt.Client, message mqtt.Message) {
		println("RemoteToLocal received")
		appName := message.Topic()[strings.LastIndex(message.Topic(), "/")+1:]

		//check if appName is in the config
		found := false
		for _, appCfg := range cfg.AppConfigs {
			if appCfg.AppName == appName {
				found = true
				break
			}
		}

		if !found {
			log.Printf("App %s not found in config", appName)
			return
		}

		server.Server.Publish("to/"+appName, message.Payload(), false, 2)
	}

	pubToRemoteTopic := cfg.RemoteBroker.BridgeTopic + "/remoteToLocal/#"
	client.Subscribe(pubToRemoteTopic, remoteToLocal, 2)

}
