// +build ignore

package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
)

// find network interface of computer mainly for Windows
func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, dev := range devices {
		fmt.Println("\nName: ", dev.Name)
		fmt.Println("Description: ", dev.Description)
		fmt.Println("Devices addresses: ", dev.Description)
		for _, addr := range dev.Addresses {
			fmt.Println("- IP address: ", addr.IP)
		}
	}
}
