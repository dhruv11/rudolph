package main

func getHelp() string {
	helpText := "I can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Help - @rudolph help"
	return helpText
}

func getContribute() string {
	contributeText := "Sorry buddy, I don't know how to do that yet, why don't you contribute to my code base? \nhttps://github.com/dhruv11/rudolph\n"
	return contributeText + getHelp()
}
