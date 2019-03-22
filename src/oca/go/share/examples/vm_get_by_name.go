package main

import (
	"fmt"
	"log"

	"github.com/OpenNebula/one/src/oca/go/src/goca"
)

func main() {
	conf := goca.NewConfig("", "", "")
	client := goca.NewClient(conf)
	controller := goca.NewController(client)

	// Get VM id by name
	id, err := controller.VMByName("vm_name")
	if err != nil {
		log.Fatal(err)
	}

	// Fetch VM informations
	vm, err := controller.VM(id).Info()
	if err != nil {
		log.Fatal(err)
	}

	// Do some stuffs on vm
	fmt.Printf("%+v\n", vm)
}
