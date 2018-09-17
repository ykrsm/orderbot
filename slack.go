package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
)

const (
	// action is used for slack attament action.

	actionStart  = "start"
	orderStart   = "orderStart"
	actionCancel = "cancel"
)

// SlackListener is
type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
}

// ListenAndResponse listens slack events and response
// particular messages. It replies by slack message button.
func (s *SlackListener) ListenAndResponse() {
	rtm := s.client.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

// handleMesageEvent handles message events.
func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	// Only response in specific channel. Ignore else.
	if ev.Channel != s.channelID {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}

	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 || m[0] != "order" {
		return fmt.Errorf("invalid message")
	}
	fmt.Printf("User: %s", ev.User)
	fmt.Printf("event:\n %+v", ev)

	// value is passed to message handler when request is approved.
	attachment := slack.Attachment{
		Text:       "Want to order something?",
		Color:      "#f9a41b",
		CallbackID: "order",
		Actions: []slack.AttachmentAction{
			{
				Name: orderStart,
				Text: "Yes!",
				Type: "button",
			},
			{
				Name:  actionCancel,
				Text:  "Cancel",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			attachment,
		},
	}

	if _, err := s.postEphemeral(ev.Channel, ev.User, "", params); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	return nil

}

func (s *SlackListener) postEphemeral(channel, user, text string, params slack.PostMessageParameters) (string, error) {
	return s.client.PostEphemeral(
		channel,
		user,
		slack.MsgOptionText(text, params.EscapeText),
		slack.MsgOptionAttachments(params.Attachments...),
		slack.MsgOptionPostMessageParameters(params),
	)
}
