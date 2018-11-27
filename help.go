package main

import (
	"math/rand"
	"strings"

	"github.com/pkg/errors"
)

func getHelp() string {
	helpText := "I can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Help - @rudolph help"
	return helpText
}

func getRating() string {
	r := [5]string{"Needs Improvement :(", "Inconsistent Performance :(", "Valued Contributor", "High Impact :)", "Exceptional Contributor :)"}
	i := rand.Intn(len(r))
	return r[i]
}

func getContribute() string {
	contributeText := "Sorry buddy, I don't know how to do that yet, why don't you contribute to my code base? \nhttps://github.com/dhruv11/rudolph\n"
	return contributeText + getHelp()
}

func wakeUp(user string, slack SlackRTMInterface) (string, error) {
	user = strings.TrimPrefix(user, "wake up <@")
	user = strings.TrimSuffix(user, ">")
	user = strings.ToUpper(user)

	_, _, c, err := slack.OpenIMChannel(user)
	if err != nil {
		return "", errors.Wrapf(err, "Could not open an IM channel to: %s", user)
	}
	slack.SendMessage(slack.NewOutgoingMessage("buddy stop napping at work, people are looking for you...", c))

	u, err := slack.GetUserInfo(user)
	if err != nil {
		return "I've just pinged them for you :)", nil
	}
	return "I've just pinged " + u.RealName + " for you :)", nil
}

func getRandomUserFromChannel(channel string, slack SlackRTMInterface) (string, error) {
	c, err := slack.GetChannelInfo(channel)
	if err != nil {
		return "", errors.Wrapf(err, "Could not retrieve channel info for: %s", channel)
	}
	i := rand.Intn(len(c.Members))

	u, err := slack.GetUserInfo(c.Members[i])
	if err != nil {
		return "", err
	}
	return u.RealName, nil
}
