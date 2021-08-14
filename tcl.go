package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gempir/go-twitch-irc/v2"
	structs "github.com/jintaokoong/go-tcl/structs"
	"github.com/jintaokoong/go-tcl/utils"
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
	client.OnConnect(func() {
		file, err := utils.GetFile("system")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		sl.SetOutput(file)
		sl.Println("connected")
		log.Println("connected")
	})
	client.OnReconnectMessage(func(message twitch.ReconnectMessage) {
		file, err := utils.GetFile("system")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		sl.SetOutput(file)
		sl.Println(message.Raw)
		log.Println(message.Raw)
	})
	client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		file, err := utils.GetFile(m.Channel)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		fl.SetOutput(file)
		un := ""
		if m.User.DisplayName == m.User.Name {
			un = m.User.Name
		} else {
			un = fmt.Sprintf("%s(%s)", m.User.DisplayName, m.User.Name)
		}
		fl.Printf("%s %s", un, m.Message)
		log.Printf("[%s] %s %s", m.Channel, un, m.Message)
	})
	client.OnClearChatMessage(func(m twitch.ClearChatMessage) {
		file, err := utils.GetFile(m.Channel)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		fl.SetOutput(file)
		fl.Println(m.Message)
		log.Println(m.Message)
	})
	client.OnNoticeMessage(func(m twitch.NoticeMessage) {
		file, err := utils.GetFile(m.Channel)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		fl.SetOutput(file)
		fl.Println(m.Message)
		log.Println(m.Message)
	})
	client.Join(config.Channels...)

	/* finally connect */
	err = client.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect()
}
