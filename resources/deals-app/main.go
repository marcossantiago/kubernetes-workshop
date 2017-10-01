package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Stores all the sweet deals
var deals []Deal

// Entrypoint of our (micro)sevice
func main() {
	fmt.Println("Deals service started...")

	initData()

	if len(deals) < 1 {
		fmt.Println("No deal. :(")
		return
	}

	http.HandleFunc("/deals", dealsHandler)
	http.HandleFunc("/healthz", healthHandler)

	http.ListenAndServe(":8080", nil)
}

// Handle deal requests
func dealsHandler(w http.ResponseWriter, r *http.Request) {
fmt.Println("request receiveed")
	idStr := r.FormValue("id")
	var id int
	var err error
	var deal Deal
	if idStr == "" {
		// Fetch random deal
		rand.Seed(time.Now().UTC().UnixNano())
		id = rand.Intn(len(deals))
		fmt.Printf("No Id passed in, fetching random: %d\n", id)
	} else {
		id, _ = strconv.Atoi(idStr)

		if err != nil {
			fmt.Println("Oops: ", err.Error())
			w.WriteHeader(500)
			return
		}
	}
	deal, err = fetchDeal(id)
	if err != nil {
		fmt.Println("Oops: ", err.Error())
		w.WriteHeader(404)
		return
	}

	err = json.NewEncoder(w).Encode(deal)
	if err != nil {
		fmt.Println("Oops: ", err.Error())
		w.WriteHeader(500)
	}
	fmt.Println("success.")
}

// Handle healthcheck
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

// Retrieve Deal by ID
func fetchDeal(id int) (Deal, error) {
	if id > len(deals) || id < 0 {
		return Deal{}, errors.New("Invalid Id")
	}
	return deals[id], nil
}

// Initialize dummy data
func initData() {
	// Load from json file
	jsonConfig, err := ioutil.ReadFile("deals.json")
	if err != nil {
		fmt.Println("Error reading data: ", err)
		return
	}
	var dealsConfig DealConfig
	err = json.Unmarshal(jsonConfig, &dealsConfig)
	if err != nil {
		fmt.Println("Error reading json config.")
		return
	}

	for i, deal := range dealsConfig.Deals {
		deals = append(deals, Deal{Id: i, Name: deal})
	}
}

type Deal struct {
	Id   int
	Name string
}

type DealConfig struct {
	Deals []string
}
