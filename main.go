package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	meetupsListID   = "5b6140b0ff2ec75df864657f"
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

	// check all external meetups and send out a reminder for any today
	resp, err := s.getMeetupReminders()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else if resp != "" {
		time.Sleep(1000 * time.Millisecond)
		s.slack.SendMessage(s.slack.NewOutgoingMessage(resp, "CBLRCPPRQ"))
	}

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
	cards, err := getCards(s.trello, listID)
	if err != nil {
		return "", errors.Wrapf(err, "Could not get card titles for list: %s", listID)
	}

	var response strings.Builder
	for _, c := range cards {
		response.WriteString(c.Name)
		response.WriteString("\n")
	}

	return response.String(), nil
}

func (s *server) getMeetupReminders() (string, error) {
	cards, err := getCards(s.trello, meetupsListID)
	if err != nil {
		return "", errors.Wrapf(err, "Could not get card titles for list: %s", meetupsListID)
	}

	var response strings.Builder
	var title = "It's your lucky day, we have a meetup later today:\n"
	response.WriteString(title)
	for _, c := range cards {
		y, m, d := c.Due.Date()

		// TODO: fix bug at month's boundary, localize date
		yy, mm, dd := time.Now().Date()
		if y == yy && m == mm && d == dd+1 {
			response.WriteString(c.Name)
			response.WriteString("\n")
		}
	}

	if response.String() == title {
		return "", nil
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

type meetup struct {
	name string
}

// TODO: replace with call to api.meetup.com, also add the date to the trello card
func (s *server) addMeetup(client *http.Client, url string) (string, error) {
	url = strings.TrimPrefix(url, "<")
	url = strings.TrimSuffix(url, ">")
	url = strings.Replace(url, "https://meetup.com", "http://api.meetup.com", -1)

	resp, err := client.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "Could not make request to %s", url)
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", errors.Wrapf(err, "Could not read request for %s", url)
	}

	var m meetup
	e := json.Unmarshal(data, &m)
	if e != nil {
		fmt.Println("couldnt deserialize" + e.Error())
	}
	fmt.Printf("got the name %s \n", m.name)

	// d := string(data)
	// st := strings.Index(d, "property=\"og:title\" content=")
	// fin := strings.Index(d[st+28:], "/>")
	// title := d[st+29:st+27+fin] + " - " + url

	// err = s.trello.CreateCard(&trello.Card{Name: title, IDList: meetupsListID}, trello.Defaults())
	// if err != nil {
	// 	return "", errors.Wrapf(err, "Could not create card with title: %s", title)
	// }

	return "looks like you just shared a meetup, I've added it to the trello board for you :)", nil
}
