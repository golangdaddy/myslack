package myslack

import (
	"fmt"
	"strings"
	//
	"github.com/nlopes/slack"
	"github.com/hokaccha/go-prettyjson"
)

type SlackEvent struct {
	APIAppID    string   `json:"api_app_id"`
	AuthedUsers []string `json:"authed_users"`
	Event       struct {
		Channel     string `json:"channel"`
		ClientMsgID string `json:"client_msg_id"`
		EventTs     string `json:"event_ts"`
		Team        string `json:"team"`
		Text        string `json:"text"`
		Ts          string `json:"ts"`
		Type        string `json:"type"`
		User        string `json:"user"`
	} `json:"event"`
	EventID   string `json:"event_id"`
	EventTime int    `json:"event_time"`
	TeamID    string `json:"team_id"`
	Token     string `json:"token"`
	Type      string `json:"type"`
	// challenge
	Challenge string `json:"challenge"`
	authedUser string
}

func (self *SlackEvent) ReplyFrom(authedUser string) {
	self.authedUser = authedUser
}

func (self *SlackEvent) Reply(text string, images ...string) error {
	return self.sendMessage(
		text,
		images...,
	)
}

func (self *SlackEvent) Replyf(formatting string, x ...interface{}) error {
	return self.sendMessage(
		fmt.Sprintf(formatting, x...),
		"",
	)
}

func (self *SlackEvent) ReplyJSON(x interface{}) error {
	b, err := prettyjson.Marshal(x)
	if err != nil {
		return err
	}
	return self.sendMessage(string(b), "")
}

func (self *SlackEvent) ReplyJSONf(formatting string, x ...interface{}) error {

	a := make([]interface{}, len(x))
	for n, _ := range x {
		b, err := prettyjson.Marshal(x)
		if err != nil {
			return err
		}
		a[n] = string(b)
	}

	return self.sendMessage(
		fmt.Sprintf(formatting, a...),
		"",
	)
}

func (self *SlackEvent) sendMessage(text string, images ...string) error {

	attachments := []slack.Attachment{}
	for _, img := range images {
		attachments = append(
			attachments,
			slack.Attachment{ImageURL: img},
		)
	}

	// https://api.slack.com/docs/message-formatting
	text = strings.Replace(text, "&", "&amp;", -1)
	text = strings.Replace(text, "<", "&lt;", -1)
	text = strings.Replace(text, ">", "&gt;", -1)

	options := []slack.MsgOption{
		slack.MsgOptionText(text, true),
	}

	if len(images) > 0 {
		options = append(
			options,
			slack.MsgOptionAttachments(attachments...),
		)
	}
	_, _, err := slack.New(self.authedUser).PostMessage(
		self.Event.Channel,
		options...,
	)
	if err != nil {
		return fmt.Errorf("FAILED TO SEND THE SLECK MESSAGE: %v", err)
	}

	return nil
}
