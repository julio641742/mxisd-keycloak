package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type KeycloakUsersArray []KeycloakSingleUserJson

type KeycloakSingleUserJson struct {
	ID         string              `json:"id"`
	Username   string              `json:"username"`
	FirstName  string              `json:"firstName"`
	LastName   string              `json:"lastName"`
	Email      string              `json:"email"`
	Enabled    bool                `json:"enabled"`
	Attributes map[string][]string `json:"attributes"`
}

func decodeKeycloakUsersArray(body io.ReadCloser) (KeycloakUsersArray, error) {
	decoder := json.NewDecoder(body)
	var kc KeycloakUsersArray
	if decoder.Decode(&kc) != nil {
		return nil, errors.New("Failed to parse json")
	}
	if len(kc) == 0 {
		return nil, errors.New("Empty user array")
	}
	return kc, nil
}

func getKeycloakToken() (string, error) {
	return loginKeycloak(userhelper, passhelper)
}

func loginKeycloak(username string, password string) (string, error) {
	token, err := config.PasswordCredentialsToken(context.Background(), username, password)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("keycloak login failed")

	}
	return token.AccessToken, nil
}

func getRequest(url string, accessToken string) (io.ReadCloser, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.New("The HTTP request failed with error")
	}
	return response.Body, nil
}

func getUserArray(token string, username string) (KeycloakSingleUserJson, error) {
	url := usersEndpoint + "?username=" + username
	body, err := getRequest(url, token)
	if err != nil {
		return KeycloakSingleUserJson{}, err
	}
	arr, err := decodeKeycloakUsersArray(body)
	defer body.Close()
	if err != nil {
		return KeycloakSingleUserJson{}, err
	}
	if arr[0].Enabled == false {
		return KeycloakSingleUserJson{}, errors.New("Account is disabled")
	}
	return arr[0], nil
}

func getUsersArray(token string) (KeycloakUsersArray, error) {
	body, err := getRequest(usersEndpoint, token)
	if err != nil {
		return nil, err
	}
	arr, err := decodeKeycloakUsersArray(body)
	defer body.Close()
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func prepareResponse(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
