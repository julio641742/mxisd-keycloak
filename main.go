package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

var keycloakURL = os.Getenv("KEYCLOAK_URL")
var realm = os.Getenv("KEYCLOAK_REALM")
var userhelper = os.Getenv("KEYCLOAK_USER")
var passhelper = os.Getenv("KEYCLOAK_PASSWORD")
var valid3pids = os.Getenv("KEYCLOAK_ATTRIBUTES_VALID_3PIDS")

var usersEndpoint = keycloakURL + "/auth/admin/realms/" + realm + "/users"

var config = oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:  keycloakURL + "/auth/realms/" + realm + "/protocol/openid-connect/auth",
		TokenURL: keycloakURL + "/auth/realms/" + realm + "/protocol/openid-connect/token",
	},
}

func main() {
	if valid3pids == "" {
		valid3pids = "email,msisdn"
	}
	if keycloakURL == "" || realm == "" || userhelper == "" || passhelper == "" ||
		os.Getenv("CLIENT_ID") == "" || os.Getenv("CLIENT_SECRET") == "" {
		fmt.Println("All required enviromental variables are not set")
		return
	}
	http.HandleFunc("/_mxisd/backend/api/v1/auth/login", Authentication)
	http.HandleFunc("/_mxisd/backend/api/v1/directory/user/search", Directory)
	http.HandleFunc("/_mxisd/backend/api/v1/identity/single", Single3PIDLookUp)
	http.HandleFunc("/_mxisd/backend/api/v1/identity/bulk", Bulk3PIDLookUp)
	http.HandleFunc("/_mxisd/backend/api/v1/profile/displayName", Profile)
	http.HandleFunc("/_mxisd/backend/api/v1/profile/threepids", Profile)
	http.HandleFunc("/_mxisd/backend/api/v1/profile/roles", Profile)

	fmt.Println("Backend is running on port 8091")
	http.ListenAndServe(":8091", nil)
}
