package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetHelp(t *testing.T) {
	expected := "I can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Recognizing a HWR behaviour - @rudolph hwr <user handle> <2 letter behaviour initial> <message> \n\tEg. @rudolph hwr @ruskin.dantra CC It was awesome when you rapped for all of us \n Help - @rudolph help"

	actual := getHelp()
	if actual != expected {
		t.Errorf("Help text was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestGetContribute(t *testing.T) {
	expected := "Sorry buddy, I don't know how to do that yet, why don't you contribute to my code base? \nhttps://github.com/dhruv11/rudolph\nI can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Recognizing a HWR behaviour - @rudolph hwr <user handle> <2 letter behaviour initial> <message> \n\tEg. @rudolph hwr @ruskin.dantra CC It was awesome when you rapped for all of us \n Help - @rudolph help"

	actual := getContribute()
	if actual != expected {
		t.Errorf("Contribute text was incorrect, got: %s, want: %s.", actual, expected)
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

func TestGetSharePrice(t *testing.T) {
	expected := "atm nzx: $11.99"

	f := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "https://www.google.co.nz/search?q=atm+nzx" {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`<span style="font-size:157%"><b>11.99</b>`)),
			}, nil
		}
		return nil, errors.New("unexpected request")
	}

	client := &http.Client{
		Transport: RoundTripFunc(f),
	}

	actual, err := getSharePrice(client, "atm nzx")

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("share price is incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestGetScheduledUpdate(t *testing.T) {
	expected := "atm nzx: $11.99\nxro asx: $11.99\n"

	f := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "https://www.google.co.nz/search?q=atm+nzx" ||
			req.URL.String() == "https://www.google.co.nz/search?q=xro+asx" {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`<span style="font-size:157%"><b>11.99</b>`)),
			}, nil
		}
		return nil, errors.New("unexpected request")
	}

	client := &http.Client{
		Transport: RoundTripFunc(f),
	}

	actual, err := getScheduledUpdate(client)

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("scheduled update is incorrect, got: %s, want: %s.", actual, expected)
	}
}

type mockClock struct {
	t time.Time
}

func (c mockClock) Now() time.Time                                { return c.t }
func (c mockClock) LoadLocation(l string) (*time.Location, error) { return time.LoadLocation(l) }

func TestShouldSendUpdate(t *testing.T) {
	tests := map[string]struct {
		input  time.Time
		output bool
	}{
		"weekday - trading hours": {
			input:  time.Date(2018, 9, 17, 0, 30, 0, 0, time.UTC),
			output: true,
		},
		"weekday - after trading hours": {
			input:  time.Date(2018, 9, 17, 14, 30, 0, 0, time.UTC),
			output: false,
		},
		"weekend": {
			input:  time.Date(2018, 9, 15, 0, 30, 0, 0, time.UTC),
			output: false,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		output := shouldSendUpdate(mockClock{t: test.input})
		assert.Equal(t, test.output, output)
	}
}

/*
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

func TestExecuteGetHelpDefault(t *testing.T) {
	expected := "helpText"

	actual, err := execute("rudolph blah", "rudolph", nil, nil, getHelpStub, nil)

	if err != nil {
		t.Errorf(err.Error())
	}
	if actual != expected {
		t.Errorf("help text is incorrect, got: %s, want: %s.", actual, expected)
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

*/
