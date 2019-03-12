package main

import (
	"math/rand"
	"strings"

	"github.com/pkg/errors"
)

func getHelp() string {
	helpText := "I can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Recognizing a HWR behaviour - @rudolph hwr <user handle> <2 letter behaviour initial> <message> \n\tEg. @rudolph hwr @ruskin.dantra CC It was awesome when you rapped for all of us \n Help - @rudolph help"
	return helpText
}

func getRating() string {
	r := [5]string{"Needs Improvement :(", "Inconsistent Performance :(", "Valued Contributor", "High Impact :)", "Exceptional Contributor :)"}
	i := rand.Intn(len(r))
	return r[i]
}

func getRisk() string {
	return "Risk is collectively owned by everyone at ASB...\nexcept for Breyten and Dhruv, they've owned enough Risk till 2021"
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

func hwr(text string, slack SlackRTMInterface) (string, error) {
	text = strings.TrimPrefix(text, "hwr <@")
	nameEnd := strings.Index(text, ">")
	name := text[0:nameEnd]
	name = strings.ToUpper(name)
	behaviours := map[string]string{
		"rr": "recognized reveller",
		"dd": "dedicated discoverer",
		"cc": "crystal clear carer",
		"pp": "punter passion",
		"ge": "gutsy evolver",
		"ra": "rapid adapter",
	}
	b := text[nameEnd+2 : nameEnd+4]
	m := text[nameEnd+5:]

	_, _, c, err := slack.OpenIMChannel(name)
	if err != nil {
		return "", errors.Wrapf(err, "Could not open an IM channel to: %s", name)
	}
	slack.SendMessage(slack.NewOutgoingMessage("Wohoo! Someone just nominated you for being a "+behaviours[b]+"!\n They said \""+m+"\"", c))

	u, err := slack.GetUserInfo(name)
	if err != nil {
		return "I've passed on your feedback anonymously! \n Good on you for being a Recognized Reveller :)", nil
	}
	return "I've passed on your feedback to " + u.RealName + " anonymously! \n Good on you for being a Recognized Reveller :)", nil
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
