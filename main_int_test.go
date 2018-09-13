package main

import (
	"testing"

	"github.com/adlio/trello"
	"github.com/dhruv11/rudolph/mocks"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/mock"
)

func TestAddInt(t *testing.T) {
	rtm := new(mocks.SlackRTMInterface)
	trelloClient := new(mocks.TrelloClient)

	srv := server{
		trello: trelloClient,
		slack:  rtm,
		done:   make(chan struct{}),
	}

	// Arrange
	incoming := make(chan slack.RTMEvent)
	rtm.On("GetIncomingEvents").Return(incoming)

	info := &slack.Info{User: &slack.UserDetails{ID: "kal"}}
	rtm.On("GetInfo").Return(info)
	rtm.On("SendMessage", mock.Anything)

	// Expectations
	rtm.On("NewOutgoingMessage", mock.MatchedBy(func(text string) bool {
		return text == "easy, your idea is in there!"
	}), mock.Anything).Return(nil)

	trelloClient.On("CreateCard", mock.MatchedBy(func(card *trello.Card) bool {
		return card.Name == "kal was here"
	}), mock.Anything).Return(nil)

	srv.start()

	msg := &slack.MessageEvent{}
	msg.Text = "<@kal> add kal was here"
	incoming <- slack.RTMEvent{
		Type: slack.TYPE_MESSAGE,
		Data: msg,
	}

	srv.stop()

	trelloClient.AssertExpectations(t)
	rtm.AssertExpectations(t)
}
