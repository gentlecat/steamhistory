// Package steam provides interface to some parts of Steam Web API.
package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	scheme = "https" // https or http
	host   = "api.steampowered.com"
)

// GetUserCount returns current number of users for a specified application.
func GetUserCount(appID int) (int, error) {
	url := scheme + "://" + host + "/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=" + strconv.Itoa(appID)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, errors.New(fmt.Sprintf("Request to %s failed (%s)!", url, resp.Status))
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return 0, err
	}

	// Types for parsing
	type response struct {
		Result       int
		Player_count int
	}
	type jason struct {
		Response response
	}
	respParsed := jason{}
	err = json.Unmarshal(body, &respParsed)
	if err != nil {
		return 0, err
	}
	return respParsed.Response.Player_count, nil
}

type App struct {
	ID   int    `json:"appid"`
	Name string `json:"name"`
}

// GetApps returns slice of all application that are available on Steam platform.
func GetApps() ([]App, error) {
	resp, err := http.Get(scheme + "://" + host + "/ISteamApps/GetAppList/v2/")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Types for parsing
	type appList struct {
		Apps []App
	}
	type jason struct {
		Applist appList
	}
	respParsed := jason{}
	err = json.Unmarshal(body, &respParsed)
	if err != nil {
		return nil, err
	}
	return respParsed.Applist.Apps, nil
}
