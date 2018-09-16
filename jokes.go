package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func getDadJoke(client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Could not make request for %s", req.URL))
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", errors.Wrapf(err, "Could not read request for %s", req.URL)
	}
	return string(data), nil
}
