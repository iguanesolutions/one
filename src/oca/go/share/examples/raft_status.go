package main

import (
	"fmt"
	"log"

	"github.com/OpenNebula/one/src/oca/go/src/goca"
)

// This example shows one way to retrieve the raft status for each server of the HA
// at end it displays the name of the leader.
func main() {
	conf := goca.NewConfig("", "", "")

	zoneName := "zone_name"

	// Create first client on the floating ip to
	// retrieve the global zone informations
	client := goca.NewClient(conf)
	controller := goca.NewController(client)

	// Retrieve the id of the zone
	id, err := controller.ZoneByName(zoneName)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve zone informations
	zone, err := controller.Zone(id).Info()
	if err != nil {
		log.Fatal(err)
	}

	// Store the endpoint of each front server of the HA
	clients := make([]*goca.Client, 0)

	// For each server of the zone which has it's own endpoint
	// we create a new client because
	for _, server := range zone.ServerPool {

		// Create a client
		conf.Endpoint = server.Endpoint
		client = goca.NewClient(conf)
		clients = append(clients, client)

		// Pass it to the controller
		controller.Client = client

		// Fetch the raft status of the server behind the endpoint
		status, err := controller.Zones().ServerRaftStatus()
		if err != nil {
			log.Fatal(err)
		}

		// Display the Raft state of the server: Leader, Follower, Candidate, Error
		state := goca.ZoneServerRaftState(status.StateRaw)
		fmt.Printf("server: %s, state: %s\n", server.Name, state.String())
	}
}
