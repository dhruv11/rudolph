// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import trello "github.com/adlio/trello"

// TrelloClient is an autogenerated mock type for the TrelloClient type
type TrelloClient struct {
	mock.Mock
}

// CreateCard provides a mock function with given fields: card, extraArgs
func (_m *TrelloClient) CreateCard(card *trello.Card, extraArgs trello.Arguments) error {
	ret := _m.Called(card, extraArgs)

	var r0 error
	if rf, ok := ret.Get(0).(func(*trello.Card, trello.Arguments) error); ok {
		r0 = rf(card, extraArgs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetList provides a mock function with given fields: listID, args
func (_m *TrelloClient) GetList(listID string, args trello.Arguments) (*trello.List, error) {
	ret := _m.Called(listID, args)

	var r0 *trello.List
	if rf, ok := ret.Get(0).(func(string, trello.Arguments) *trello.List); ok {
		r0 = rf(listID, args)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*trello.List)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, trello.Arguments) error); ok {
		r1 = rf(listID, args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
