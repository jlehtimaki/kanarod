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
	db.accessToken = os.Getenv("ACCESS_TOKEN")
	if db.accessToken == "" {
		log.Fatal("could not find ACCESS_TOKEN environment variable")
	}
	db.apiPath = os.Getenv("API_PATH")
	if db.apiPath == "" {
		log.Infof("using the default API Path %s", defaultApiPath)
		db.apiPath = defaultApiPath
		return
	}
	log.Infof("using API_PATH: %s", db.apiPath)
}

func (d *discordBot) getNextOpponent(team string) (string, error) {
	apiPath := fmt.Sprintf("%s/match/next/%s", d.apiPath, team)
	data, err := d.restGet(apiPath)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (d *discordBot) getTeam(team string) (string, error) {
	apiPath := fmt.Sprintf("%s/team/%s", d.apiPath, team)
	data, err := d.restGet(apiPath)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (d *discordBot) restGet(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", d.accessToken)
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
