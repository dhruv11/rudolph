package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type passenger struct {
	Name    string
	Address string
}

func getPassengers(client *http.Client) (string, error) {
	resp, err := client.Get("http://prod.j22cbjqtiv.us-east-1.elasticbeanstalk.com/passengers")
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Could not make request to get passengers"))
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", errors.Wrapf(err, "Could not read request to get passengers")
	}

	var p []passenger
	e := json.Unmarshal(data, &p)
	if e != nil {
		return "", errors.Wrapf(e, "Could not deserialise request to get passengers")
	}

	var r = "Your choices are:\n"
	for i := range p {
		r = r + p[i].Name + " from " + p[i].Address + "\n"
	}
	return r, nil
}
