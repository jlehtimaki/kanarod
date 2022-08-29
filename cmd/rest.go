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

type Rest struct {
	AccessToken string
	ApiPath     string
}

func initRest() Rest {
	rest := Rest{}
	rest.AccessToken = os.Getenv("ACCESS_TOKEN")
	if rest.AccessToken == "" {
		log.Fatal("could not find ACCESS_TOKEN environment variable")
	}
	rest.ApiPath = os.Getenv("API_PATH")
	if rest.ApiPath == "" {
		log.Infof("using the default API Path %s", defaultApiPath)
		rest.ApiPath = defaultApiPath
		return rest
	}
	log.Infof("using API_PATH: %s", rest.ApiPath)
	return rest
}

func (r *Rest) getNextOpponent(team string) (string, error) {
	apiPath := fmt.Sprintf("%s/match/next/%s", r.ApiPath, team)
	data, err := r.restGet(apiPath)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (r *Rest) getTeam(team string) (string, error) {
	apiPath := fmt.Sprintf("%s/team/%s", r.ApiPath, team)
	data, err := r.restGet(apiPath)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (r *Rest) restGet(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", r.AccessToken)
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

func (r *Rest) getScheduled(team string) (string, error) {
	apiPath := fmt.Sprintf("%s/match/scheduled/%s", r.ApiPath, team)
	data, err := r.restGet(apiPath)
	if err != nil {
		return "", err
	}
	return data, err
}
