package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	structs "github.com/jintaokoong/go-tcl/structs"
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
		log.Println("connected")
	})
	client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		ct := time.Now()
		fn := fmt.Sprintf("%s_%s.log", m.Channel, ct.Format("20060201"))
		file, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		fl := log.New(file, "", log.Ldate|log.Ltime)
		fl.Printf("%s(%s) %s", m.User.DisplayName, m.User.Name, m.Message)
		log.Printf("[%s] %s(%s) %s", m.Channel, m.User.DisplayName, m.User.Name, m.Message)
	})
	client.Join(config.Channels...)

	/* finally connect */
	err = client.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect()
}
