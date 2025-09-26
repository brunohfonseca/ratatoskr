package notifications

import (
	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
)

func SendSlackMsg(cfg *config.AppConfig, attachmentBody slack.Attachment) error {
	client := slack.New(cfg.Alerts.Slack.Token)
	attachment := attachmentBody

	_, _, err := client.PostMessage(cfg.Alerts.Slack.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		log.Error().Msgf("Failed to send message to Slack: %v", err)
		return err
	}
	return nil
}
