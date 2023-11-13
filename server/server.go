package server

import (
	"bytes"
	"log"

	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/thomaskhub/muecke/client"
	"github.com/thomaskhub/muecke/utils"
)

var done chan bool

type MochiHooks struct {
	mochi.HookBase
}

func (h *MochiHooks) OnStarted() {
	println("mochi started")
	done <- true
}
func (h *MochiHooks) OnSubscribe(cl *mochi.Client, pk packets.Packet) packets.Packet {
	return pk
}

func (h *MochiHooks) OnConnect(cl *mochi.Client, pk packets.Packet) error {
	return nil
}

func (h *MochiHooks) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mochi.OnStarted,
	}, []byte{b})
}

type MochiServer struct {
	Server *mochi.Server
	Client *client.MqttClient
}

func (srv *MochiServer) StartServer(cfg *utils.MueckeConfig, mqttClient *client.MqttClient) {
	done = make(chan bool, 1)

	srv.Server = mochi.New(&mochi.Options{
		InlineClient: true,
	})

	srv.Client = mqttClient

	srv.Server.AddHook(new(auth.AllowHook), nil)

	go func() {
		tcp := listeners.NewTCP("t1", "127.0.0.1:1883", nil)
		srv.Server.AddListener(tcp)
		srv.Server.AddHook(new(MochiHooks), map[string]any{})
		err := srv.Server.Serve()
		if err != nil {
			log.Fatal(err)
		}

	}()

	//wait for done = true so that we know mochi server is started
	<-done

	//only after broker is running we can subscribe
	for i, appCfg := range cfg.AppConfigs {

		remoteTopic := cfg.RemoteBroker.BridgeTopic + "/localToRemote/" + appCfg.AppName
		topic := "from/" + appCfg.AppName

		srv.Server.Subscribe(topic, i, func(cl *mochi.Client, sub packets.Subscription, pk packets.Packet) {
			mqttClient.Publish(remoteTopic, string(pk.Payload), 2)
		})
	}

}
