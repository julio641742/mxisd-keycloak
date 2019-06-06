package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthRequest struct {
	Auth struct {
		Mxid      string `json:"mxid"`
		Localpart string `json:"localpart"`
		Domain    string `json:"domain"`
		Password  string `json:"password"`
	} `json:"auth"`
}

type AuthResponse struct {
	Auth AuthResponseBody `json:"auth"`
}

type AuthResponseBody struct {
	Success bool         `json:"success"`
	ID      *LookUpID    `json:"id,omitempty"`
	Profile *ProfileBody `json:"profile,omitempty"`
}

func Authentication(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	var req AuthRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		fmt.Println("Invalid json request")
		return
	}
	res := AuthResponse{
		Auth: AuthResponseBody{
			Success: false,
		},
	}
	_, err := loginKeycloak(req.Auth.Localpart, req.Auth.Password)
	if err != nil {
		fmt.Println(err)
	} else {
		profile, err := buildProfile(req.Auth.Localpart)
		if err != nil {
			fmt.Println(err)
		} else {
			res = AuthResponse{
				Auth: AuthResponseBody{
					Success: true,
					ID: &LookUpID{
						Type:  "localpart",
						Value: req.Auth.Localpart,
					},
					Profile: profile,
				},
			}
		}
	}
	prepareResponse(w, res)

}
