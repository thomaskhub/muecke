package test

import (
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thomaskhub/muecke/bridge"
	"github.com/thomaskhub/muecke/client"
	"github.com/thomaskhub/muecke/utils"
)

func TestBridge(t *testing.T) {
	cfg := utils.ParseConfig("../config.yaml")

	bridge.StartBridge(cfg)

	// Create two clients
	client1 := client.NewMqttClientWithConfig(
		cfg.RemoteBroker.Broker,
		cfg.RemoteBroker.ClientId+"manfred",
		cfg.RemoteBroker.Username,
		cfg.RemoteBroker.Password,
	)

	client2 := client.NewMqttClientWithConfig(
		cfg.RemoteBroker.Broker,
		cfg.RemoteBroker.ClientId+"2"+"manfred",
		cfg.RemoteBroker.Username,
		cfg.RemoteBroker.Password,
	)

	// Create channels to receive messages
	ch1 := make(chan string)
	ch2 := make(chan string)

	client1.Connect()
	client2.Connect()

	// Subscribe client1 to remote-bridge/localToRemote/app1
	client1.Subscribe(cfg.RemoteBroker.BridgeTopic+"/localToRemote/App1", func(client mqtt.Client, message mqtt.Message) {
		ch1 <- string(message.Payload())
	}, 2)

	// // Subscribe client2 to to/app1
	client2.Subscribe("to/App1", func(client mqtt.Client, message mqtt.Message) {
		ch2 <- string(message.Payload())
	}, 2)

	client1.Publish(cfg.RemoteBroker.BridgeTopic+"/remoteToLocal/App1", "testing123", 2)

	// Wait for message on client2 and publish it to from/app1
	go func() {
		msg := <-ch2
		client2.Publish("from/App1", msg, 2)
	}()

	// Wait for message on client1
	select {
	case msg := <-ch1:
		if msg != "testing123" {
			t.Errorf("Expected 'testing123', got '%s'", msg)
		}
	case <-time.After(time.Second * 5):
		t.Error("Timeout waiting for message")
	}

}
