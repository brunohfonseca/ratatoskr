package notifications

import (
	"strconv"

	"github.com/brunohfonseca/ratatoskr/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func SendTelegramMsg(cfg *config.AppConfig, message string) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Alerts.Telegram.BotToken)
	if err != nil {
		log.Error().Msgf("Failed to create bot: %v", err)
		return err
	}
	chatID, err := strconv.ParseInt(cfg.Alerts.Telegram.ChatID, 10, 64)
	if err != nil {
		log.Error().Msgf("Failed to parse chat ID: %v", err)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"

	_, err = bot.Send(msg)
	if err != nil {
		log.Error().Msgf("Failed to send message to Telegram: %v", err)
		return err
	}
	return nil
}
