package main

import (
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
func getCards(client TrelloClient, listID string) ([]*trello.Card, error) {
	list, err := client.GetList(listID, trello.Defaults())
	if err != nil {
		return nil, err
	}

	cards, err := list.GetCards(trello.Defaults())
	if err != nil {
		return nil, err
	}

	return cards, nil
}
