package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/gorilla/websocket"
)

type Message struct {
	Os      string `json:"os"`
	Counter uint64 `json:"counter"`
}

var (
	windowsCounter uint64
	androidCounter uint64
	iosCounter     uint64

	messages [3]Message

	// ethOS map[net.HardwareAddr]string
	ethOS map[string]string = map[string]string{}

	// variable for packet capture
	snapshotLen int32 = 1024
	promiscuous bool  = true
	err         error
	timeout     time.Duration = 30 * time.Second
	handle      *pcap.Handle

	ethernetLayer  gopacket.Layer
	ipLayer        gopacket.Layer
	dhcpLayer      gopacket.Layer
	ethernetPacket *layers.Ethernet
	ipPacket       *layers.IPv4
	dhcpPacket     *layers.DHCPv4

	srcMac  string
	os_name []string
	host_os string
)

func sendCounter(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrade, _ := upgrader.Upgrade(w, r, nil)

	messages[0].Os = "Windows"
	messages[1].Os = "Android"
	messages[2].Os = "iOS"

	for {
		// send number of packet every 1 second
		time.Sleep(1 * time.Second)

		// messages = messages[:0] // clear slice
		messages[0].Counter = windowsCounter
		messages[1].Counter = androidCounter
		messages[2].Counter = iosCounter

		jsonBytes, err := json.Marshal(messages)
		if err != nil {
			log.Fatal(err)
		}
		jsonStr := string(jsonBytes)

		upgrade.WriteJSON(jsonStr)

		iosCounter = 0
		androidCounter = 0
		windowsCounter = 0
	}
}

func capturePacket(device string, fd *os.File) {
	log.Printf("Start capturing packets")

	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packets := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packets.Packets() {
		// fmt.Println(packet)
		countPacket(packet, fd)
	}
}

func countPacket(packet gopacket.Packet, fd *os.File) {
	ethernetLayer = packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer == nil {
		return
	}
	ethernetPacket, _ = ethernetLayer.(*layers.Ethernet)
	srcMac = (ethernetPacket.SrcMAC).String()
	if os, ok := ethOS[srcMac]; ok {
		switch os {
		case "Windows":
			windowsCounter++
		case "Android":
			androidCounter++
		case "iOS":
			iosCounter++
		}
	}

	dhcpLayer = packet.Layer(layers.LayerTypeDHCPv4)
	if dhcpLayer != nil {
		if _, ok := ethOS[srcMac]; ok {
			// already know os of host computer
			return
		}
		dhcpPacket, _ = dhcpLayer.(*layers.DHCPv4)
		dhcpFingerprinting(srcMac, dhcpPacket, fd)
	}
}

func dhcpFingerprinting(srcMac string, dhcpPacket *layers.DHCPv4, fd *os.File) {
	os_name = os_name[:0]
	for _, option := range dhcpPacket.Options {
		if option.Type == layers.DHCPOptHostname {
			os_name = append(os_name, string(option.Data))
		} else if option.Type == layers.DHCPOptClassID {
			os_name = append(os_name, string(option.Data))
		}
	}

	if len(os_name) == 0 {
		return
	}

	for _, os := range os_name {
		switch {
		case strings.Contains(os, "MSFT"):
			host_os = "Windows"
		case strings.Contains(os, "android"):
			host_os = "Android"
		case strings.Contains(os, "iPhone"):
			host_os = "iOS"
		case strings.Contains(os, "iphone"):
			host_os = "iOS"
		case strings.Contains(os, "MBP"):
			host_os = "OSX"
		}
	}

	if len(host_os) == 0 {
		return
	}

	fmt.Println("os_name: ", os_name)
	fmt.Printf("---- srcMac: [%s], os: [%s] ----\n", string(srcMac), host_os)
	ethOS[string(srcMac)] = host_os

	// write log to file
	str := fmt.Sprintf("srcMac: [%s], os: [%s]\n", string(srcMac), host_os)
	fd.WriteString(str)
	host_os = ""
}
