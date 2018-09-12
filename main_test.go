package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
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
	expectedListID   string
}

func (client testTrelloClient) GetCardTitles(listID string) ([]string, error) {
	if client.unhappyPath {
		return nil, errors.New("unhappy")
	}

	if listID != client.expectedListID {
		return nil, errors.New("list id is incorrect")
	}

	cards := []string{"card1", "card2"}
	return cards, nil
}

func TestGetListItems(t *testing.T) {
	expected := "card1\ncard2\n"

	actual, err := getListItems("123", testTrelloClient{expectedListID: "123"})

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("list items are incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestGetListItemsUnhappy(t *testing.T) {
	_, err := getListItems("123", testTrelloClient{unhappyPath: true})

	if err == nil {
		t.Errorf("expected an error")
	}
}

func (client testTrelloClient) CreateCard(title string, listID string) error {
	if client.unhappyPath {
		return errors.New("unhappy")
	}

	if title != client.expectedCardName {
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

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGetDadJoke(t *testing.T) {
	expected := "haha"

	f := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "https://icanhazdadjoke.com/" && req.Header.Get("Accept") == "text/plain" {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`haha`)),
			}, nil
		}
		return nil, errors.New("unexpected request")
	}

	client := &http.Client{
		Transport: RoundTripFunc(f),
	}

	actual, err := getDadJoke(client)

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("joke is incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestGetDadJokeUnhappy(t *testing.T) {
	f := func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("unhappy path")
	}

	client := &http.Client{
		Transport: RoundTripFunc(f),
	}

	_, err := getDadJoke(client)

	if err == nil {
		t.Errorf("expected an error")
	}
}

func getHelpStub() string {
	return "helpText"
}

func getContributeStub() string {
	return "contributeText"
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

func TestExecuteGetHelpDefault(t *testing.T) {
	expected := "contributeText"

	actual, err := execute("rudolph blah", "rudolph", nil, nil, nil, nil, getContributeStub)

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("contribute text is incorrect, got: %s, want: %s.", actual, expected)
	}
}

func getDadJokeStub(client *http.Client) (string, error) {
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

func getListItemsStub(listID string, client trelloClientAdapter) (string, error) {
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

func addIdeaStub(title string, client trelloClientAdapter) (string, error) {
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
