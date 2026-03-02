package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func NewBot(cfg Config, log *logrus.Logger) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(cfg.TG.BotToken)
	if err != nil {
		log.Fatalf("Failed to start telegram bot: %+v", err)
	}

	return bot
}

func Start(bot *tgbotapi.BotAPI, cfg Config, log *logrus.Logger) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go func(updates tgbotapi.UpdatesChannel) {
		for update := range updates {
			if update.Message != nil && update.Message.IsCommand() {
				if update.Message.Command() == "start" {

					text := "<b>👋 Добро пожаловать в ИСПО TOOLS!</b>\n\nНажмите кнопку ниже, чтобы открыть расписание."
					webappURL := cfg.TG.WebUrl

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
					msg.ParseMode = "HTML"
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.InlineKeyboardButton{
						Text: "⚡️ Перейти в приложение",
						URL:  &webappURL,
					}))

					if _, err := bot.Send(msg); err != nil {
						log.Errorf("Failed to send message via bot: %+v", err)
					}
				}
			}
		}
	}(updates)
}
