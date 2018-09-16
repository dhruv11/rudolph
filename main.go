package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/adlio/trello"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

const (
	ideasListID     = "5b613db79ea6a782ac173a48"
	scheduledListID = "5b613dbfd923da512f85263b"
)

func main() {
	srv := newServer()
	srv.start()
	<-srv.done
}

type server struct {
	trello TrelloClient
	slack  SlackRTMInterface
	done   chan struct{}
}

func newServer() server {
	appKey := os.Getenv("TRELLO_KEY")
	trello := os.Getenv("TRELLO_TOKEN")
	slack := os.Getenv("SLACK_TOKEN")

	return server{
		trello: newTrelloClient(appKey, trello),
		slack:  newSlackRTM(slack),
		done:   make(chan struct{}),
	}
}

func (s *server) start() {
	fmt.Println("Starting")

	go func(done chan struct{}, s *server) {
		for {
			select {
			// channel operator, await the async goroutine
			case <-s.done:
				fmt.Println("Stopping")
				return

			case event := <-s.slack.GetIncomingEvents():
				switch msg := event.Data.(type) {
				case *slack.ConnectedEvent:
					fmt.Println("Connection counter:", msg.ConnectionCount)

				case *slack.MessageEvent:
					info := s.slack.GetInfo()
					prefix := fmt.Sprintf("<@%s> ", info.User.ID)
					resp, err := s.processMessage(msg, info, prefix)
					if err != nil {
						fmt.Printf("Error: %s\n", err)
					}
					if msg.User != info.User.ID && strings.HasPrefix(msg.Text, prefix) {
						s.slack.SendMessage(s.slack.NewOutgoingMessage(resp, msg.Channel))
					}

				case *slack.RTMError:
					fmt.Printf("Error: %s\n", msg.Error())

				case *slack.LatencyReport:
					// send scheduled updates to me
					if shouldSendUpdate() {
						r, err := getScheduledUpdate(&http.Client{})
						if err != nil {
							fmt.Printf("Error: %s\n", err)
						} else {
							s.slack.SendMessage(s.slack.NewOutgoingMessage(r, "DCKGBPU10"))
						}
					}

				case *slack.InvalidAuthEvent:
					fmt.Printf("Invalid credentials")
					close(done)

				default:
					// do nothing
				}
			}
		}
	}(s.done, s)
}

func (s *server) stop() {
	close(s.done)
}

func (s *server) processMessage(msg *slack.MessageEvent, info *slack.Info, prefix string) (string, error) {
	text := strings.TrimPrefix(msg.Text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	if strings.HasSuffix(text, "scheduled") {
		return s.getListItems(scheduledListID)
	} else if strings.HasSuffix(text, "ideas") {
		return s.getListItems(ideasListID)
	} else if strings.HasPrefix(text, "add") {
		return s.addIdea(text)
	} else if strings.HasPrefix(text, "price") {
		return getSharePrice(&http.Client{}, text)
	} else if text == "make me laugh" {
		return getDadJoke(&http.Client{})
	} else if text == "help" {
		return getHelp(), nil
	}
	return getContribute(), nil
}

func (s *server) getListItems(listID string) (string, error) {
	titles, err := getCardTitles(s.trello, listID)
	if err != nil {
		return "", errors.Wrapf(err, "Could not get card titles for list: %s", listID)
	}

	var response strings.Builder
	for _, t := range titles {
		response.WriteString(t)
		response.WriteString("\n")
	}

	return response.String(), nil
}

func (s *server) addIdea(title string) (string, error) {
	title = strings.TrimPrefix(title, "add")
	title = strings.TrimSpace(title)

	err := s.trello.CreateCard(&trello.Card{Name: title, IDList: ideasListID}, trello.Defaults())
	if err != nil {
		return "", errors.Wrapf(err, "Could not create card with title: %s", title)
	}

	return "easy, your idea is in there!", nil
}

func getHelp() string {
	helpText := "I can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Help - @rudolph help"
	return helpText
}

func getContribute() string {
	contributeText := "Sorry buddy, I don't know how to do that yet, why don't you contribute to my code base? \nhttps://github.com/dhruv11/rudolph\n"
	return contributeText + getHelp()
}

func shouldSendUpdate() bool {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		fmt.Println("Could not find timezone")
		return false
	}
	now := time.Now().In(loc)

	h := []int{22, 0, 2, 4}
	if now.Weekday() < 5 && contains(h, now.Hour()) &&
		now.Minute() == 30 && now.Second() < 30 {
		return true
	}
	return false
}

func contains(arr []int, h int) bool {
	for _, a := range arr {
		if a == h {
			return true
		}
	}
	return false
}
