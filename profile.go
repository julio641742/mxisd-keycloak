package main

import (
    "fmt"
    "net/http"
    "encoding/json"
)

type ProfileRequest struct {
    Mxid string `json:"mxid"`
    Localpart string `json:"localpart"`
    Domain string `json:"domain"`
}

type ProfileResponse struct {
    Profile ProfileResponseBody `json:"profile"`
}

type ProfileResponseBody struct {
    DisplayName string `json:"display_name,omitempty"`
    ThreePIDS []Profile3PIDS `json:"threepids,omitempty"`
    Roles []string `json:"roles,omitempty"`
}

func Profile(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        decoder := json.NewDecoder(r.Body)
        var req ProfileRequest
        err := decoder.Decode(&req)
        if err != nil {
            fmt.Println("Invalid json request")
            return
        }
        var res ProfileResponse
        token, err := getKeycloakToken()
        if err != nil {
            fmt.Println(err)
        } else {
            user, err := getUserArray(token, req.Localpart)
            if err != nil {
                fmt.Println(err)
            } else {
                roles, err := getProfileRoles(token, user.Id)
                if err != nil {
                    fmt.Println(err)
                } else {
                    res = ProfileResponse {
                        Profile: ProfileResponseBody {
                            DisplayName: getDisplayName(user),
                            ThreePIDS: getProfile3PIDS(user),
                            Roles: roles,
                        },
                    }
                }
            }
        }
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(res)
    }
}