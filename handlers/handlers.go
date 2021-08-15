package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/jintaokoong/go-tcl/structs"
	"github.com/jintaokoong/go-tcl/utils"
)

func HandleConnect(sl *log.Logger) func() {
	return func() {
		file, err := utils.GetFile("system")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		sl.SetOutput(file)
		sl.Println("connected")
		log.Println("connected")
	}
}

func HandleReconnectMessage(sl *log.Logger) func(message twitch.ReconnectMessage) {
	return func(message twitch.ReconnectMessage) {
		file, err := utils.GetFile("system")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		sl.SetOutput(file)
		sl.Println(message.Raw)
		log.Println(message.Raw)
	}
}

func HandlePrivateMessage(sl *log.Logger, fl *log.Logger, config structs.Config) func(m twitch.PrivateMessage) {
	return func(m twitch.PrivateMessage) {
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

		sysLog := []string{fmt.Sprintf("#%s", m.Channel)}
		fileLog := []string{}

		for _, role := range config.Roles {
			if m.User.Badges[role] == 0 {
				continue
			}
			pr := strings.ToUpper(utils.GetFirstN(role, 3))
			sysLog = append(sysLog, fmt.Sprintf("[%s]", pr))
			fileLog = append(fileLog, fmt.Sprintf("[%s]", pr))
		}

		sysLog = append(sysLog, un, m.Message)
		fileLog = append(fileLog, un, m.Message)

		log.Print(strings.Join(sysLog, " "))
		fl.Print(strings.Join(fileLog, " "))
	}
}

func HandleClearChatMessage(sl *log.Logger, fl *log.Logger) func(m twitch.ClearChatMessage) {
	return func(m twitch.ClearChatMessage) {
		file, err := utils.GetFile(m.Channel)
		if err != nil {
			sl.Println(err)
			log.Panic(err)
		}
		defer file.Close()
		fl.SetOutput(file)
		fl.Println(m.Message)
		log.Println(fmt.Sprintf("#%s", m.Channel), m.Message)
	}
}

func HandleNoticeMessage(sl *log.Logger, fl *log.Logger) func(m twitch.NoticeMessage) {
	return func(m twitch.NoticeMessage) {
		file, err := utils.GetFile(m.Channel)
		if err != nil {
			sl.Println(err)
			log.Panic(err)
		}
		defer file.Close()
		fl.SetOutput(file)
		fl.Println(m.Message)
		log.Println(fmt.Sprintf("#%s", m.Channel), m.Message)
	}
}
