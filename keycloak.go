package main

import (
    "io"
    "fmt"
    "errors"
    "context"
    "strings"
    "net/http"
    "encoding/json"
)

type KeycloakUsersArray []KeycloakSingleUserJson

type KeycloakSingleUserJson struct {
    Id string `json:"id"`
    Username string `json:"username"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
    Email string `json:"email"`
    Enabled bool `json:"enabled"`
    Attributes map[string][]string `json:"attributes"`
}

type Profile3PIDS struct {
    Medium string `json:"medium,omitempty"`
    Address string `json:"address,omitempty"`
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
        return "", errors.New("Keycloak login failed!")
        
    }
    return token.AccessToken, nil
}

func getRequest(url string, accessToken string) (io.ReadCloser, error) {
    client := &http.Client{}
    request, _ := http.NewRequest("GET", url, nil)
    request.Header.Set("Authorization", "Bearer " + accessToken)
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

func getDisplayName(user KeycloakSingleUserJson) (string) {
    name :=  strings.TrimSpace(user.FirstName + " " + user.LastName)
    if name == "" {
        name = user.Username
    }
    return name
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

func findUserBy3PID(medium string, address string, users KeycloakUsersArray, token string) (*LookUpResponseBodyId) {
    if users == nil {
        users, _ = getUsersArray(token)
        if users == nil {
            return nil
        }
    }
    // lookup.Type can either be localpart or mxid (only localpart is supported for now)
    for _, user := range users {
        if medium == "email" && user.Email == address {
            return &LookUpResponseBodyId {
                Type: "localpart",
                Value: user.Username,
            }
        }
        for key, value := range user.Attributes {
            if medium == key && address == value[0]  {
                return &LookUpResponseBodyId {
                    Type: "localpart",
                    Value: user.Username,
                }
            }
        }
    }
    return nil
}

func getProfile3PIDS(user KeycloakSingleUserJson) ([]Profile3PIDS) {
    var tpids []Profile3PIDS
    if user.Email != "" {
        tpids = append(tpids, Profile3PIDS {
            Medium: "email",
            Address: user.Email,
        })
    }
    for key, value := range user.Attributes {
        tpids = append(tpids, Profile3PIDS {
            Medium: key,
            Address: value[0],
        })
    }
    return tpids
}

func getProfileRoles(token string, id string) ([]string, error) {
    url := usersEndpoint + "/" + id + "/groups"
    body, err := getRequest(url, token)
    if err != nil {
        return nil, err
    }
    decoder := json.NewDecoder(body)
    var kc []struct {
        Name string `json:"name"`
    }
    if decoder.Decode(&kc) != nil {
        return nil, errors.New("Failed to parse json")
    }
    defer body.Close()
    roles := make([]string, len(kc))
    for i, v := range kc {
        roles[i] = v.Name
    }
    return roles, nil
}
