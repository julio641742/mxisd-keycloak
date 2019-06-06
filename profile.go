package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type ProfileRequest struct {
	Mxid      string `json:"mxid"`
	Localpart string `json:"localpart"`
	Domain    string `json:"domain"`
}

type ProfileResponse struct {
	Profile ProfileBody `json:"profile"`
}

type Profile3PIDS struct {
	Medium  string `json:"medium,omitempty"`
	Address string `json:"address,omitempty"`
}

type ProfileBody struct {
	DisplayName string         `json:"display_name,omitempty"`
	AvatarURL   string         `json:"avatar_url,omitempty"`
	ThreePIDS   []Profile3PIDS `json:"threepids,omitempty"`
	Roles       []string       `json:"roles,omitempty"`
}

func getDisplayName(user KeycloakSingleUserJson) string {
	name := strings.TrimSpace(user.FirstName + " " + user.LastName)
	if name == "" {
		name = user.Username
	}
	return name
}

func getAvatarURL(attributes map[string][]string) string {
	if _, value := attributes["avatar_url"]; value {
		return attributes["avatar_url"][0]
	}
	return ""
}

func getProfile3PIDS(user KeycloakSingleUserJson) []Profile3PIDS {
	var tpids []Profile3PIDS
	if user.Email != "" {
		tpids = append(tpids, Profile3PIDS{
			Medium:  "email",
			Address: user.Email,
		})
	}
	for key, value := range user.Attributes {
		if strings.Contains(valid3pids, key) {
			tpids = append(tpids, Profile3PIDS{
				Medium:  key,
				Address: value[0],
			})
		}
	}
	return tpids
}

// Keycloak groups are used as roles
func getProfileRoles(token string, id string) ([]string, error) {
	url := usersEndpoint + "/" + id + "/groups"
	body, err := getRequest(url, token)
	if err != nil {
		return nil, err
	}
	var kc []struct {
		Name string `json:"name"`
	}
	if json.NewDecoder(body).Decode(&kc) != nil {
		return nil, errors.New("Failed to parse json")
	}
	defer body.Close()
	roles := make([]string, len(kc))
	for i, v := range kc {
		roles[i] = v.Name
	}
	return roles, nil
}

func buildProfile(username string) (*ProfileBody, error) {
	token, err := getKeycloakToken()
	if err != nil {
		return nil, err
	}
	user, err := getUserArray(token, username)
	if err != nil {
		return nil, err
	}
	roles, err := getProfileRoles(token, user.ID)
	if err != nil {
		return nil, err
	}
	return &ProfileBody{
		DisplayName: getDisplayName(user),
		AvatarURL:   getAvatarURL(user.Attributes),
		ThreePIDS:   getProfile3PIDS(user),
		Roles:       roles,
	}, nil
}

func Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	var req ProfileRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		fmt.Println("Invalid json request")
		return
	}
	var res ProfileResponse
	profile, err := buildProfile(req.Localpart)
	if err != nil {
		fmt.Println(err)
	} else {
		res.Profile = *profile
	}
	prepareResponse(w, res)

}
