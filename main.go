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

const (
	ideasListID     = "5b613db79ea6a782ac173a48"
	scheduledListID = "5b613dbfd923da512f85263b"
	meetupsListId   = "5b6140b0ff2ec75df864657f"
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

					resp, err := s.processMessage(msg, info, prefix, s.slack)
					if err != nil {
						fmt.Printf("Error: %s\n", err)
					}
					// TODO: Move the prefix check to processMessage
					if msg.User != info.User.ID && (strings.HasPrefix(msg.Text, prefix) || strings.HasPrefix(msg.Text, "<https://www.meetup.com/")) {
						s.slack.SendMessage(s.slack.NewOutgoingMessage(resp, msg.Channel))
					}

				case *slack.RTMError:
					fmt.Printf("Error: %s\n", msg.Error())

				case *slack.LatencyReport:
					// send scheduled updates to me
					if shouldSendUpdate(realClock{}) {
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

func (s *server) processMessage(msg *slack.MessageEvent, info *slack.Info, prefix string, slack SlackRTMInterface) (string, error) {
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
	} else if strings.HasPrefix(text, "wake up") {
		return wakeUp(text, slack)
	} else if strings.HasPrefix(text, "<https://www.meetup.com/") {
		return s.addMeetup(&http.Client{}, text)
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

func (s *server) addMeetup(client *http.Client, url string) (string, error) {
	url = strings.TrimPrefix(url, "<")
	url = strings.TrimSuffix(url, ">")

	resp, err := client.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "Could not make request to %s", url)
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", errors.Wrapf(err, "Could not read request for %s", url)
	}

	d := string(data)
	st := strings.Index(d, "property=\"og:title\" content=")
	fin := strings.Index(d[st+28:], "/>")
	title := d[st+29:st+27+fin] + " - " + url

	err = s.trello.CreateCard(&trello.Card{Name: title, IDList: meetupsListId}, trello.Defaults())
	if err != nil {
		return "", errors.Wrapf(err, "Could not create card with title: %s", title)
	}

	return "looks like you just shared a meetup, I've added it to the trello board for you :)", nil
}
