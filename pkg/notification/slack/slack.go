package slack

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nlopes/slack"

	// "github.com/keel-hq/keel/constants"
	"github.com/rusenask/cloudstore/pkg/notification"
	"github.com/rusenask/cloudstore/types"

	log "github.com/sirupsen/logrus"
)

const (
	EnvSlackToken            = "SLACK_TOKEN"
	EnvSlackBotName          = "SLACK_BOT_NAME"
	EnvSlackChannels         = "SLACK_CHANNELS"
	EnvSlackApprovalsChannel = "SLACK_APPROVALS_CHANNEL"
)

const timeout = 5 * time.Second

type sender struct {
	slackClient *slack.Client
	channels    []string
	botName     string
}

func init() {
	notification.RegisterSender("slack", &sender{})
}

func (s *sender) Configure(config *notification.Config) (bool, error) {
	var token string
	// Get configuration
	if os.Getenv(EnvSlackToken) != "" {
		token = os.Getenv(EnvSlackToken)
	} else {
		return false, nil
	}
	if os.Getenv(EnvSlackBotName) != "" {
		s.botName = os.Getenv(EnvSlackBotName)
	} else {
		s.botName = "keel"
	}

	if os.Getenv(EnvSlackChannels) != "" {
		channels := os.Getenv(EnvSlackChannels)
		s.channels = strings.Split(channels, ",")
	} else {
		s.channels = []string{"general"}
	}

	s.slackClient = slack.New(token)

	log.WithFields(log.Fields{
		"name":     "slack",
		"channels": s.channels,
	}).Info("extension.notification.slack: sender configured")

	return true, nil
}

func (s *sender) Send(event types.EventNotification) error {
	params := slack.NewPostMessageParameters()
	params.Username = s.botName
	params.IconURL = "https://cdn4.iconfinder.com/data/icons/robots-ultra-colour-collection/60/Robots_-_Ultra_Color_-_015_-_Cylon-512.png"

	params.Attachments = []slack.Attachment{
		slack.Attachment{
			Fallback: event.Message,
			Color:    event.Level.Color(),
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "New upload!",
					Value: event.Message,
					Short: false,
				},
			},
			Footer: "keel.sh",
			Ts:     json.Number(strconv.Itoa(int(event.CreatedAt.Unix()))),
		},
	}

	chans := s.channels
	if len(event.Channels) > 0 {
		chans = event.Channels
	}

	for _, channel := range chans {
		_, _, err := s.slackClient.PostMessage(channel, "", params)
		if err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"channel": channel,
			}).Error("extension.notification.slack: failed to send notification")
		}
	}
	return nil
}
