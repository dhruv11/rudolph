package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/adlio/trello"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

const ideasListID = "5b613db79ea6a782ac173a48"
const scheduledListID = "5b613dbfd923da512f85263b"

var client trelloClientAdapter

func main() {

	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	api.SetDebug(true)

	rtm := api.NewRTM()
	// goroutine, async exec
	go rtm.ManageConnection()

Loop:
	for {
		select {
		// channel operator, await the async goroutine
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					response, err := execute(ev.Text, prefix, getListItems, addIdea, getHelp, getDadJoke)
					if err != nil {
						fmt.Printf("Error: %s\n", err)
					} else {
						rtm.SendMessage(rtm.NewOutgoingMessage(response, ev.Channel))
					}
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//do nothing
			}
		}
	}
}

func execute(text string, prefix string, getListItems listGetter,
	addIdea cardCreator, getHelp helpGetter, getDadJoke jokeGetter) (string, error) {
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	if strings.HasSuffix(text, "scheduled") {
		return getListItems(scheduledListID, getTrelloClient())
	} else if strings.HasSuffix(text, "ideas") {
		return getListItems(ideasListID, getTrelloClient())
	} else if strings.HasPrefix(text, "add") {
		return addIdea(text, getTrelloClient())
	} else if text == "make me laugh" {
		return getDadJoke()
	}
	return getHelp(), nil
}

type listGetter func(listID string, client trelloClientAdapter) (string, error)

func getListItems(listID string, client trelloClientAdapter) (string, error) {
	titles, err := client.GetCardTitles(listID)
	if err != nil {
		return "", err
	}

	var response strings.Builder
	for _, t := range titles {
		response.WriteString(t)
		response.WriteString("\n")
	}

	return response.String(), nil
}

type cardCreator func(title string, client trelloClientAdapter) (string, error)

func addIdea(title string, client trelloClientAdapter) (string, error) {
	title = strings.TrimPrefix(title, "add")
	title = strings.TrimSpace(title)

	err := client.CreateCard(title, ideasListID)
	if err != nil {
		return "", err
	}

	return "easy, your idea is in there!", nil
}

type helpGetter func() string

func getHelp() string {
	helpText := "I can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Help - @rudolph help \n Feature request - @dhruv <request>"

	return helpText
}

type jokeGetter func() (string, error)

func getDadJoke() (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "text/plain")

	response, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Could not make request for %s", req.URL))
	}

	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Could not read request for %s", req.URL))
	}
	return string(data), nil
}

type trelloClientAdapter interface {
	CreateCard(title string, listID string) error
	GetCardTitles(listID string) ([]string, error)
}

type trelloClient struct {
	client *trello.Client
}

func (client trelloClient) CreateCard(title string, listID string) error {
	err := client.client.CreateCard(&trello.Card{Name: title, IDList: listID}, trello.Defaults())
	if err != nil {
		return err
	}
	return nil
}

func (client trelloClient) GetCardTitles(listID string) ([]string, error) {
	list, err := client.client.GetList(listID, trello.Defaults())
	if err != nil {
		return nil, err
	}

	cards, err := list.GetCards(trello.Defaults())
	if err != nil {
		return nil, err
	}

	t := []string{}
	for _, c := range cards {
		t = append(t, c.Name)
	}
	return t, nil
}

func getTrelloClient() trelloClientAdapter {
	if client == nil {
		appKey := os.Getenv("TRELLO_KEY")
		token := os.Getenv("TRELLO_TOKEN")

		client = trelloClient{client: trello.NewClient(appKey, token)}
	}
	return client
}
