package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	defaultApiPath = "http://localhost:8080"
)

func initRest(db *discordBot) {
	db.apiPath = os.Getenv("API_PATH")
	if db.apiPath == "" {
		log.Infof("using the default API Path %s", defaultApiPath)
		db.apiPath = defaultApiPath
	}
}

func (d *discordBot) getNextOpponent(team string) (string, error) {
	apiPath := fmt.Sprintf("%s/match/next/%s", d.apiPath, team)
	data, err := restGet(apiPath)
	if err != nil {
		return "", err
	}
	return data, nil
}

func restGet(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf(resp.Status)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}
	//Convert the body to type string
	return string(body), nil
}
