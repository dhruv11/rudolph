package main

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func getDadJoke() (string, error) {
	u := "https://icanhazdadjoke.com/"
	resp, err := http.Get(u)
	if err != nil {
		return "", errors.Wrapf(err, "Could not make request to %s", u)
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", errors.Wrapf(err, "Could not read request for %s", u)
	}
	return string(data), nil
}
