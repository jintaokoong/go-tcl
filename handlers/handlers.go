package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

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

func HandlePrivateMessageV2(config structs.Config) func(m twitch.PrivateMessage) {
	return func(m twitch.PrivateMessage) {
		mongoClient, err := utils.CreateClient(config)
		if err != nil {
			log.Panic("create client error", err)
		}

		ctx, cancel := utils.CreateDatabaseContext()
		defer cancel()

		err = mongoClient.Connect(ctx)
		if err != nil {
			log.Panic("connection error", err)
		}
		defer mongoClient.Disconnect(ctx)

		database := mongoClient.Database(config.Database.Name)
		collectionName := fmt.Sprintf("%s%s", config.Database.Collection, time.Now().Local().Format("200601"))
		entryCollection := database.Collection(collectionName)

		un := ""
		if m.User.DisplayName == m.User.Name {
			un = m.User.Name
		} else {
			un = fmt.Sprintf("%s(%s)", m.User.DisplayName, m.User.Name)
		}
		sysLog := []string{fmt.Sprintf("#%s", m.Channel)}
		roles := []string{}
		for _, role := range config.Roles {
			if m.User.Badges[role] == 0 {
				continue
			}
			pr := strings.ToUpper(utils.GetFirstN(role, 3))
			sysLog = append(sysLog, fmt.Sprintf("[%s]", pr))
			roles = append(roles, role)
		}

		entry := structs.Entry{
			DisplayName:     m.User.DisplayName,
			UserID:          m.User.Name,
			Channel:         m.Channel,
			Message:         m.Message,
			Roles:           roles,
			CreatedDatetime: time.Now(),
		}
		_, err = entryCollection.InsertOne(ctx, entry)
		if err != nil {
			log.Panic(err)
		}

		sysLog = append(sysLog, un, m.Message)
		log.Print(strings.Join(sysLog, " "))
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
		log.Println(fmt.Sprintf("#%s", m.Channel), m.Raw)
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
