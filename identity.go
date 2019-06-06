package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Single3PIDLookUpRequest struct {
	Lookup Profile3PIDS `json:"lookup"`
}

type SingleLookUpResponse struct {
	Lookup *LookUpResponseBody `json:"lookup,omitempty"`
}

type LookUpID struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type LookUpResponseBody struct {
	Medium  string    `json:"medium,omitempty"`
	Address string    `json:"address,omitempty"`
	ID      *LookUpID `json:"id,omitempty"`
}

func findUserBy3PID(medium string, address string, users KeycloakUsersArray, token string) *LookUpID {
	if users == nil {
		users, _ = getUsersArray(token)
		if users == nil {
			return nil
		}
	}
	for _, user := range users {
		if medium == "email" && user.Email == address {
			return &LookUpID{
				Type:  "localpart",
				Value: user.Username,
			}
		}
		for key, value := range user.Attributes {
			if strings.Contains(valid3pids, key) && medium == key && address == value[0] {
				return &LookUpID{
					Type:  "localpart",
					Value: user.Username,
				}
			}
		}
	}
	return nil
}

func Single3PIDLookUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	var req Single3PIDLookUpRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		fmt.Println("Invalid json request")
		return
	}
	var res SingleLookUpResponse
	token, err := getKeycloakToken()
	if err != nil {
		fmt.Println(err)
	} else {
		found := findUserBy3PID(req.Lookup.Medium, req.Lookup.Address, nil, token)
		if found != nil {
			res = SingleLookUpResponse{
				Lookup: &LookUpResponseBody{
					Medium:  req.Lookup.Medium,
					Address: req.Lookup.Address,
					ID:      found,
				},
			}
		}
	}
	prepareResponse(w, res)

}

type Bulk3PIDLookUpRequest struct {
	Lookup []Profile3PIDS `json:"lookup"`
}

type BulkLookUpResponse struct {
	Lookup []LookUpResponseBody `json:"lookup"`
}

func Bulk3PIDLookUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	var req Bulk3PIDLookUpRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		fmt.Println("Invalid json request")
		return
	}
	var res = BulkLookUpResponse{
		Lookup: make([]LookUpResponseBody, 0),
	}
	token, err := getKeycloakToken()
	if err != nil {
		fmt.Println(err)
	} else {
		users, err := getUsersArray(token)
		if err != nil {
			fmt.Println(err)
		} else {
			var results []LookUpResponseBody
			for _, tpid := range req.Lookup {
				found := findUserBy3PID(tpid.Medium, tpid.Address, users, "")
				if found != nil {
					results = append(results, LookUpResponseBody{
						Medium:  tpid.Medium,
						Address: tpid.Address,
						ID:      found,
					},
					)
				}
			}
			if results != nil {
				res.Lookup = results
			}
		}
	}
	prepareResponse(w, res)

}
