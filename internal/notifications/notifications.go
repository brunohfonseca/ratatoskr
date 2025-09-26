package notifications

import (
	"github.com/bhfonseca/ratatoskr/internal/config"
	"github.com/slack-go/slack"
)

func SendAlert(cfg *config.AppConfig, message, logType string) error {
	err := SendTelegramMsg(cfg, message)
	if err != nil {
		return err
	}
	err = SendSlackMsg(cfg, configureSlackAttachment(message, logType))
	if err != nil {
		return err
	}
	return nil
}

func configureSlackAttachment(message, logType string) slack.Attachment {
	attachment := slack.Attachment{
		Color: "#36a64f", // Default color
		Text:  message,
	}

	switch logType {
	case "error":
		attachment.Color = "#ff0000" // Red for errors
	case "warning":
		attachment.Color = "#ffa500" // Orange for warnings
	case "info":
		attachment.Color = "#36a64f" // Green for info
	default:
		attachment.Color = "#cccccc" // Grey for other types
	}

	return attachment
}
