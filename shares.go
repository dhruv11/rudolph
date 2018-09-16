package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

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

type clock interface {
	Now() time.Time
	LoadLocation(l string) (*time.Location, error)
}

type realClock struct{}

func (realClock) Now() time.Time                                { return time.Now() }
func (realClock) LoadLocation(l string) (*time.Location, error) { return time.LoadLocation(l) }

func shouldSendUpdate(clock clock) bool {
	loc, err := clock.LoadLocation("UTC")
	if err != nil {
		fmt.Println("Could not find timezone")
		return false
	}
	now := clock.Now().In(loc)

	h := []int{22, 0, 2, 4}
	if now.Weekday() < 5 && contains(h, now.Hour()) &&
		now.Minute() == 30 && now.Second() < 30 {
		return true
	}
	return false
}

func contains(arr []int, h int) bool {
	for _, a := range arr {
		if a == h {
			return true
		}
	}
	return false
}
