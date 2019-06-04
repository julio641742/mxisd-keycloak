package main

import (
    "fmt"
    "net/http"
    "encoding/json"
)

type AuthRequest struct {
    Auth struct {
        Mxid string `json:"mxid"`
        Localpart string `json:"localpart"`
        Domain string `json:"domain"`
        Password string `json:"password"`
    } `json:"auth"`
}

type AuthResponse struct {
    Auth AuthResponseBody `json:"auth"`
}

type AuthResponseBody struct {
    Success bool `json:"success"`
    Id *AuthResponseBodyId `json:"id,omitempty"`
    Profile *AuthResponseBodyProfile `json:"profile,omitempty"`
}

type AuthResponseBodyId struct {
    Type string `json:"type,omitempty"`
    Value string `json:"value,omitempty"`
}

type AuthResponseBodyProfile struct {
    DisplayName string `json:"display_name,omitempty"`
    ThreePIDS   []Profile3PIDS `json:"three_pids,omitempty"`
}

func Authentication(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        decoder := json.NewDecoder(r.Body)
        var req AuthRequest
        err := decoder.Decode(&req)
        if err != nil {
            fmt.Println("Invalid json request")
            return
        }
        res := AuthResponse {
            Auth: AuthResponseBody {
                Success: false,
            },
        }
        _, err = loginKeycloak(req.Auth.Localpart, req.Auth.Password)
        if err != nil {
            fmt.Println(err)
        } else {
            token, err := getKeycloakToken()
            if err != nil {
                fmt.Println(err)
            } else {
                user, err := getUserArray(token, req.Auth.Localpart)
                if err != nil {
                    fmt.Println(err)
                } else {
                    res = AuthResponse {
                        Auth: AuthResponseBody {
                            Success: true,
                            Id: &AuthResponseBodyId {
                                Type: "localpart",
                                Value: req.Auth.Localpart,
                            },
                            Profile: &AuthResponseBodyProfile {
                                DisplayName: getDisplayName(user),
                                ThreePIDS: getProfile3PIDS(user),
                            },
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