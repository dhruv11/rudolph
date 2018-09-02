package main

import (
	"errors"
	"testing"

	"github.com/adlio/trello"
)

func TestGetHelpText(t *testing.T) {
	expected := "I can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Help - @rudolph help \n Feature request - @dhruv <request>"

	actual := getHelp()
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

func getHelpStub() string {
	return "helpText"
}

func TestExecuteGetHelp(t *testing.T) {
	expected := "helpText"

	actual, err := execute("rudolph HELP", "rudolph", nil, nil, getHelpStub, nil)

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("help text is incorrect, got: %s, want: %s.", actual, expected)
	}
}

func getDadJokeStub() (string, error) {
	return "joke", nil
}

func TestExecuteGetJoke(t *testing.T) {
	expected := "joke"

	actual, err := execute("rudolph make me laugh", "rudolph", nil, nil, nil, getDadJokeStub)

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("joke is incorrect, got: %s, want: %s.", actual, expected)
	}
}

func getListItemsStub(listID string) (string, error) {
	return "card1\ncard2", nil
}

func TestExecuteGetList(t *testing.T) {
	expected := "card1\ncard2"

	actual, err := execute("rudolph ideas", "rudolph", getListItemsStub, nil, nil, nil)

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("list is incorrect, got: %s, want: %s.", actual, expected)
	}
}

func addIdeaStub(title string, client trelloClient) (string, error) {
	return "done", nil
}

func TestExecuteAddIdea(t *testing.T) {
	expected := "done"

	actual, err := execute("rudolph add blah", "rudolph", nil, addIdeaStub, nil, nil)

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("add idea response is incorrect, got: %s, want: %s.", actual, expected)
	}
}
