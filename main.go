package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os/exec"
)

func main() {
	viper.SetConfigName("dev")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(viper.GetString("telegram.token"))
	if err != nil {
		logrus.Panic(err)
	}

	bot.Debug = viper.GetBool("general.debug")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.From.ID != 1418203928 {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Reloaded.")

		switch update.Message.Command() {
		case "reload":
			down := exec.Command("wg-quick", "down", "/etc/wireguard/wg0.conf")
			if err := down.Run(); err != nil {
				logrus.Error(err)
				continue
			}

			up := exec.Command("wg-quick", "up", "/etc/wireguard/wg0.conf")
			if err := up.Run(); err != nil {
				logrus.Error(err)
				continue
			}

			if _, err := bot.Send(msg); err != nil {
				logrus.Error(err)
				continue
			}
		}
	}
}
