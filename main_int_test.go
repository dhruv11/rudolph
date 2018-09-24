package main

import (
	"testing"
	"time"

	"github.com/adlio/trello"
	"github.com/dhruv11/rudolph/mocks"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
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

	trelloClient.On("GetList", mock.Anything, mock.Anything).Return(&trello.List{}, errors.New("throwing so we can skip this bit"))

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

	time.Sleep(100 * time.Millisecond)
	trelloClient.AssertExpectations(t)
	rtm.AssertExpectations(t)
}

func TestHelpInt(t *testing.T) {
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

	info := &slack.Info{User: &slack.UserDetails{ID: "newbie"}}
	rtm.On("GetInfo").Return(info)
	rtm.On("SendMessage", mock.Anything)

	// Expectations
	rtm.On("NewOutgoingMessage", mock.MatchedBy(func(text string) bool {
		return text == "I can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Help - @rudolph help"
	}), mock.Anything).Return(nil)

	trelloClient.On("GetList", mock.Anything, mock.Anything).Return(&trello.List{}, errors.New("throwing so we can skip this bit"))

	srv.start()

	msg := &slack.MessageEvent{}
	msg.Text = "<@newbie> help"
	incoming <- slack.RTMEvent{
		Type: slack.TYPE_MESSAGE,
		Data: msg,
	}

	srv.stop()

	time.Sleep(100 * time.Millisecond)
	rtm.AssertExpectations(t)
}

func TestContributeInt(t *testing.T) {
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

	info := &slack.Info{User: &slack.UserDetails{ID: "newbie"}}
	rtm.On("GetInfo").Return(info)
	rtm.On("SendMessage", mock.Anything)

	// Expectations
	rtm.On("NewOutgoingMessage", mock.MatchedBy(func(text string) bool {
		return text == "Sorry buddy, I don't know how to do that yet, why don't you contribute to my code base? \nhttps://github.com/dhruv11/rudolph\nI can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Help - @rudolph help"
	}), mock.Anything).Return(nil)

	trelloClient.On("GetList", mock.Anything, mock.Anything).Return(&trello.List{}, errors.New("throwing so we can skip this bit"))

	srv.start()

	msg := &slack.MessageEvent{}
	msg.Text = "<@newbie> blah"
	incoming <- slack.RTMEvent{
		Type: slack.TYPE_MESSAGE,
		Data: msg,
	}

	srv.stop()

	time.Sleep(100 * time.Millisecond)
	rtm.AssertExpectations(t)
}
