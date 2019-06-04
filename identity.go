package main

import (
    "fmt"
    "net/http"
    "encoding/json"
)

type Single3PIDLookUpRequest struct {
    Lookup Profile3PIDS `json:"lookup"`
}

type SingleLookUpResponse struct {
    Lookup *LookUpResponseBody `json:"lookup,omitempty"`
}

type LookUpResponseBody struct {
    Medium string `json:"medium,omitempty"`
    Address string `json:"address,omitempty"`
    Id *LookUpResponseBodyId `json:"id,omitempty"`
}

type LookUpResponseBodyId struct {
    Type string `json:"type,omitempty"`
    Value string `json:"value,omitempty"`
}

func Single3PIDLookUp(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        decoder := json.NewDecoder(r.Body)
        var req Single3PIDLookUpRequest
        err := decoder.Decode(&req)
        if err != nil {
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
                res = SingleLookUpResponse {
                    Lookup: &LookUpResponseBody {
                        Medium: req.Lookup.Medium,
                        Address: req.Lookup.Address,
                        Id: found,
                    },
                }
            }
        }
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(res)
    }
}

type Bulk3PIDLookUpRequest struct {
    Lookup []Profile3PIDS `json:"lookup"`
}

type BulkLookUpResponse struct {
    Lookup []LookUpResponseBody `json:"lookup"`
}

func Bulk3PIDLookUp(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        decoder := json.NewDecoder(r.Body)
        var req Bulk3PIDLookUpRequest
        err := decoder.Decode(&req)
        if err != nil {
            fmt.Println("Invalid json request")
            return
        }
        var res = BulkLookUpResponse {
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
                        results = append(results, LookUpResponseBody {
                            Medium: tpid.Medium,
                            Address: tpid.Address,
                            Id: found,
                        },
                        )
                    }
                }
                if results != nil {
                    res.Lookup = results
                }
            }
        }
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(res)
    }
}