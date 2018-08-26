package main

import (
	"errors"
	"testing"

	"github.com/adlio/trello"
)

func TestGetHelpText(t *testing.T) {
	expected := "I can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Help - @rudolph help \n Feature request - @dhruv <request>"

	actual := GetHelpText()
	if actual != expected {
		t.Errorf("Help text was incorrect, got: %s, want: %s.", actual, expected)
	}
}

type testTrelloClient struct {
}

func (testTrelloClient) CreateCard(card *trello.Card, extraArgs trello.Arguments) error {
	if card.Name != "testing" {
		return errors.New("card name is incorrect")
	}
	return nil
}

func TestAddIdea(t *testing.T) {
	expected := "easy, your idea is in there!"

	actual, err := addIdea("testing", testTrelloClient{})

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("List name was incorrect, got: %s, want: %s.", actual, expected)
	}
}
