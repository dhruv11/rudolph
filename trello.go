package main

import (
	"time"

	"github.com/adlio/trello"
)

// TrelloClient - for mocking trello client
type TrelloClient interface {
	CreateCard(card *trello.Card, extraArgs trello.Arguments) error
	GetList(listID string, args trello.Arguments) (list *trello.List, err error)
}

func newTrelloClient(appKey, token string) TrelloClient {
	return trello.NewClient(appKey, token)
}

// Could create a Trello struct to put this on, same as with slack
func getCardTitles(client TrelloClient, listID string) ([]string, error) {
	list, err := client.GetList(listID, trello.Defaults())
	if err != nil {
		return nil, err
	}

	cards, err := list.GetCards(trello.Defaults())
	if err != nil {
		return nil, err
	}

	titles := []string{}
	for _, c := range cards {
		titles = append(titles, c.Name)
	}
	return titles, nil
}

// TODO: refactor this to leverage the function above
func getCardsDueToday(client TrelloClient, listID string) ([]string, error) {
	list, err := client.GetList(listID, trello.Defaults())
	if err != nil {
		return nil, err
	}

	cards, err := list.GetCards(trello.Defaults())
	if err != nil {
		return nil, err
	}

	titles := []string{}
	for _, c := range cards {
		y, m, d := c.Due.Date()

		// TODO: fix bug at month's boundary, localize date
		yy, mm, dd := time.Now().Date()
		if y == yy && m == mm && d == dd+1 {
			titles = append(titles, c.Name)
		}
	}
	return titles, nil
}
