package main

import "github.com/nlopes/slack"

// SlackRTMInterface - for mocking slack client
type SlackRTMInterface interface {
	GetInfo() *slack.Info
	NewOutgoingMessage(text string, channelID string, options ...slack.RTMsgOption) *slack.OutgoingMessage
	SendMessage(msg *slack.OutgoingMessage)
	GetIncomingEvents() chan slack.RTMEvent
	GetUserInfo(user string) (*slack.User, error)
	OpenIMChannel(user string) (bool, bool, string, error)
	GetChannelInfo(channelID string) (*slack.Channel, error)
}

type slackRTM struct {
	rtm *slack.RTM
}

func newSlackRTM(token string) *slackRTM {
	c := slack.New(token)
	c.SetDebug(true)

	rtm := c.NewRTM()
	// goroutine, async exec
	go rtm.ManageConnection()

	return &slackRTM{
		rtm: rtm,
	}
}

func (s *slackRTM) OpenIMChannel(user string) (bool, bool, string, error) {
	return s.rtm.OpenIMChannel(user)
}

func (s *slackRTM) GetChannelInfo(channelID string) (*slack.Channel, error) {
	return s.rtm.GetChannelInfo(channelID)
}

func (s *slackRTM) GetUserInfo(user string) (*slack.User, error) {
	return s.rtm.GetUserInfo(user)
}

func (s *slackRTM) GetInfo() *slack.Info {
	return s.rtm.GetInfo()
}

func (s *slackRTM) GetIncomingEvents() chan slack.RTMEvent {
	return s.rtm.IncomingEvents
}

func (s *slackRTM) SendMessage(msg *slack.OutgoingMessage) {
	s.rtm.SendMessage(msg)
}

func (s *slackRTM) NewOutgoingMessage(text string, channelID string, options ...slack.RTMsgOption) *slack.OutgoingMessage {
	return s.rtm.NewOutgoingMessage(text, channelID, options...)
}
