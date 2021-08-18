package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/jintaokoong/go-tcl/handlers"
	"github.com/jintaokoong/go-tcl/structs"
	"gopkg.in/yaml.v2"
)

func main() {
	/* capture interrupt */
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigc
		log.Panic("exit")
	}()

	fl := log.New(nil, "", log.Ldate|log.Ltime)
	sl := log.New(nil, "", log.Ldate|log.Ltime)

	/* read config */
	log.Println("initializing")
	configFile, err := os.Open("config.yml")
	if err != nil {
		log.Panic("missing config file!")
	}
	defer configFile.Close()

	var config structs.Config
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Panic("config file invalid format")
	}

	log.Printf("joining %s", config.Channels)

	/* configure client */
	client := twitch.NewAnonymousClient()
	client.OnConnect(handlers.HandleConnect(sl))
	client.OnReconnectMessage(handlers.HandleReconnectMessage(sl))
	client.OnPrivateMessage(handlers.HandlePrivateMessageV2(config))
	client.OnClearChatMessage(handlers.HandleClearChatMessage(sl, fl))
	client.OnNoticeMessage(handlers.HandleNoticeMessage(sl, fl))
	client.Join(config.Channels...)

	/* finally connect */
	err = client.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect()
}
