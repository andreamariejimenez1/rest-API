package main

import (
	"encoding/json"
	"log" //log erros
	"net/http"

	"github.com/gorilla/mux"
)

// Creating a configuration struct (MODEL)
type Network struct {
	TheType   string `json:"TYPE"`
	Bootproto string `json:"BOOTPROTO"`
	Name      string `json:"NAME"`
	Device    string `json:"DEVICE"`
	Onboot    string `json:"ONBOOT"`
	Ipaddr    string `json:"IPADDR"`
	Prefix    string `json:"PREFIX"`
}

// Initiliazing "networks" variable of Network Struct -- variable lenght array
var networks []Network

//Get all networks
func getAllNetworks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(networks) // to write JSON to the server
}

//Get single network
func getNetwork(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Read name parameter

	// Iterate through all networks and find the one that matches the name
	for _, item := range networks {
		if item.Name == params["NAME"] { // if item name is eqaul to the NAME parameter in the url
			json.NewEncoder(w).Encode(item) // print out the single network
			return
		}
	}
	json.NewEncoder(w).Encode(&Network{})
}

//Create a new network
func createNetwork(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var network Network // create a variable "network" and set it to Network struct
	_ = json.NewDecoder(r.Body).Decode(&network)
	networks = append(networks, network) // append to global variable networks
	json.NewEncoder(w).Encode(network)   // giving us a response
}

//Update a network
func updateNetwork(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range networks { // Checking if the network exsist
		if item.Name == params["NAME"] { // by comparing the name and NAME parameter
			networks = append(networks[:index], networks[index+1:]...) // slice it out
			var network Network
			_ = json.NewDecoder(r.Body).Decode(&network)
			networks = append(networks, network)

			json.NewEncoder(w).Encode(network) // giving us a response
			return
		}
	}
	json.NewEncoder(w).Encode(networks)
}

//Delete
func deleteNetwork(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range networks {
		if item.Name == params["NAME"] {
			networks = append(networks[:index], networks[index+1:]...) //slice it out
			break
		}
	}
	// json.NewEncoder(w).Encode(networks)
	json.NewEncoder(w).Encode("Deleted")
}

func main() {
	// Initializing router
	router := mux.NewRouter()

	// Mock Data
	// TODO -- implement database
	// networks = append(networks, Network{TheType: "Ethernet", Bootproto: "static", Name: "eth0", Device: "eth0", Onboot: "yes", Ipaddr: "192.168.1.1", Prefix: "24"})

	// Route Handlers - to establish endpoints for the api
	router.HandleFunc("/sysconfig/network-scripts/ifcfg", getAllNetworks).Methods("GET")
	router.HandleFunc("/sysconfig/network-scripts/ifcfg-{NAME}", getNetwork).Methods("GET")
	router.HandleFunc("/sysconfig/network-scripts/ifcfg", createNetwork).Methods("POST")
	router.HandleFunc("/sysconfig/network-scripts/ifcfg-{NAME}", updateNetwork).Methods("PUT")
	router.HandleFunc("/sysconfig/network-scripts/ifcfg-{NAME}", deleteNetwork).Methods("DELETE")

	// To run the server
	log.Fatal(http.ListenAndServe(":8000", router))
}
