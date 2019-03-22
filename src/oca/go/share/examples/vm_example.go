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

	// Build a string template. (No XML-RPC call done)
	// To make a VM from an existing OpenNebula template,
	// use template "Instantiate" method instead
	tpl := goca.NewTemplateBuilder()
	tpl.AddValue("name", "this-is-a-vm")
	tpl.AddValue("cpu", 1)
	tpl.AddValue("vcpu", "2")
	tpl.AddValue("memory", "64")

	// The disk ID should exist to make this example work
	vec := tpl.NewVector("disk")
	vec.AddValue("image_id", "119")
	vec.AddValue("dev_prefix", "vd")

	// Create the VM from the template
	vmID, err := controller.VMs().Create(tpl.String(), false)
	if err != nil {
		log.Fatal(err)
	}

	// Keep pointer on VM controller to achieve a list of actions
	vmCtrl := controller.VM(vmID)

	// Fetch informations of the created VM
	vm, err := vmCtrl.Info()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", vm)

	// Poweroff the VM
	err = vmCtrl.PoweroffHard()
	if err != nil {
		log.Fatal(err)
	}

}
