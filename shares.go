package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func getScheduledUpdate(client *http.Client) (string, error) {
	shares := []string{"atm nzx", "xro asx"}

	var res strings.Builder
	for _, share := range shares {
		r, err := getSharePrice(client, share)
		if err != nil {
			continue
		}
		res.WriteString(r)
		res.WriteString("\n")
	}

	return res.String(), nil
}

func getSharePrice(client *http.Client, symbol string) (string, error) {
	symbol = strings.TrimPrefix(symbol, "price")
	symbol = strings.TrimSpace(symbol)

	u := fmt.Sprintf("https://www.google.co.nz/search?q=%s", url.QueryEscape(symbol))
	resp, err := client.Get(u)
	if err != nil {
		return "", errors.Wrapf(err, "Could not make request to %s", u)
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", errors.Wrapf(err, "Could not read request for %s", u)
	}

	d := string(data)
	span := strings.Index(d, "<span style=\"font-size:157%\"><b>")
	f := strings.Index(d[span+32:], "</b>")

	return symbol + ": $" + d[span+32:span+32+f], nil
}
