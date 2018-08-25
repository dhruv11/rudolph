package main

import "testing"

func TestGetHelpText(t *testing.T) {
	expected := "I can help you with: \n Fetching ideas - @rudolph ideas \n Fetching scheduled talks - @rudolph scheduled \n Adding an idea: @rudolph add <talk title> \n Dad joke - @rudolph make me laugh \n Help - @rudolph help \n Feature request - @dhruv <request>"

	actual := GetHelpText()
	if actual != expected {
		t.Errorf("Help text was incorrect, got: %s, want: %s.", actual, expected)
	}
}
