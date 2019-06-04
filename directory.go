package main

import (
    "fmt"
    "strings"
    "net/http"
    "encoding/json"
)

type DirectoryRequest struct {
    By string `json:"by"`
    SearchTerm string `json:"search_term"`
}

type DirectoryResponse struct {
    Limited bool `json:"limited"`
    Results []DirectoryResponseResults `json:"results"`
}

type DirectoryResponseResults struct {
    AvatarURL string `json:"avatar_url,omitempty"`
    DisplayName string `json:"display_name,omitempty"`
    UserID string `json:"user_id,omitempty"`
}

func existsBy3PID(user KeycloakSingleUserJson, term string) bool {
    foundByEmail := strings.Contains(user.Email, term)
    foundBy3PID := false
    for _, value := range user.Attributes {
        foundBy3PID = strings.Contains(value[0], term)
        if foundBy3PID {
            break
        }
    }
    return foundByEmail || foundBy3PID
}

func Directory(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        decoder := json.NewDecoder(r.Body)
        var req DirectoryRequest
        err := decoder.Decode(&req)
        if err != nil {
            fmt.Println("Invalid json request")
            return
        }
        var res = DirectoryResponse {
            Limited: false,
            Results: make([]DirectoryResponseResults, 0),
        }
        if req.By == "name" || req.By == "threepid" {
            token, err := getKeycloakToken()
            if err != nil {
                fmt.Println(err)
            } else {
                users, err := getUsersArray(token)
                if err != nil {
                    fmt.Println(err)
                } else {
                    var results []DirectoryResponseResults
                    for _, user := range users {
                        byname := req.By == "name" && (strings.Contains(user.Username, req.SearchTerm) || strings.Contains(user.Username, req.SearchTerm))
                        bytpid := req.By == "threepid" && existsBy3PID(user, req.SearchTerm)
                        if byname || bytpid {
                            results = append(results, DirectoryResponseResults {
                                AvatarURL: "",
                                DisplayName: getDisplayName(user),
                                UserID: user.Username,
                            },
                            )
                        }
                    }
                    if results != nil {
                        res.Results = results
                    }
                }
            }
        }
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(res)
    }
}