package main

import (
	"errors"
	"testing"

	"github.com/adlio/trello"
)

func TestGetHelpText(t *testing.T) {
	expected := "I can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Help - @rudolph help \n Feature request - @dhruv <request>"

	actual := getHelpText()
	if actual != expected {
		t.Errorf("Help text was incorrect, got: %s, want: %s.", actual, expected)
	}
}

type testTrelloClient struct {
	unhappyPath      bool
	expectedCardName string
}

func (client testTrelloClient) CreateCard(card *trello.Card, extraArgs trello.Arguments) error {
	if client.unhappyPath {
		return errors.New("unhappy")
	}

	if card.Name != client.expectedCardName {
		return errors.New("card name is incorrect")
	}
	return nil
}

func TestAddIdea(t *testing.T) {
	expected := "easy, your idea is in there!"

	actual, err := addIdea("add testing", testTrelloClient{expectedCardName: "testing"})

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("list name is incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestAddIdeaUnhappy(t *testing.T) {
	_, err := addIdea("add testing", testTrelloClient{unhappyPath: true})

	if err == nil {
		t.Errorf("expected an error")
	}
}
