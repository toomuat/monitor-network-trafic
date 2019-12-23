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

// Message : packet counter of each OS
type Message struct {
	Os      string `json:"os"`
	Counter uint64 `json:"counter"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	windowsCounter uint64
	androidCounter uint64
	iosCounter     uint64

	messages = []Message{
		{Os: "Windows"},
		{Os: "Android"},
		{Os: "iOS"},
	}

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

	srcMac string
	osName []string
	hostOS string

	sendNum   uint64
	clientNum uint64
)

func handleWS(w http.ResponseWriter, r *http.Request) {
	clientNum++
	fmt.Printf("new client [%d] joined\n", clientNum)
	clientId := clientNum

	// cannot detect leave of client
	// right now, this function keep working even if I close tab of blowser
	// following defer program doesn't work
	defer func() {
		fmt.Printf("client [%d] left\n", clientNum)
		clientNum--
	}()

	upgrade, _ := upgrader.Upgrade(w, r, nil)

	sendCounter(upgrade, clientId)
}

func sendCounter(upgrade *websocket.Conn, clientId uint64) {
	for {
		// send number of packet every 1 second
		time.Sleep(1 * time.Second)

		messages[0].Counter = windowsCounter
		messages[1].Counter = androidCounter
		messages[2].Counter = iosCounter

		jsonBytes, err := json.Marshal(messages)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(clientId)

		upgrade.WriteJSON(string(jsonBytes))
		sendNum++

		// need to be initialized if send to the all clients
		if sendNum == clientNum {
			iosCounter = 0
			androidCounter = 0
			windowsCounter = 0
			sendNum = 0
		}
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
	osName = osName[:0]
	for _, option := range dhcpPacket.Options {
		if option.Type == layers.DHCPOptHostname {
			osName = append(osName, string(option.Data))
		} else if option.Type == layers.DHCPOptClassID {
			osName = append(osName, string(option.Data))
		}
	}

	if len(osName) == 0 {
		return
	}

	for _, os := range osName {
		switch {
		case strings.Contains(os, "MSFT"):
			hostOS = "Windows"
		case strings.Contains(os, "android"):
			hostOS = "Android"
		case strings.Contains(os, "iPhone"):
			hostOS = "iOS"
		case strings.Contains(os, "iphone"):
			hostOS = "iOS"
		case strings.Contains(os, "MBP"):
			hostOS = "OSX"
		}
	}

	if len(hostOS) == 0 {
		return
	}

	fmt.Println("osName: ", osName)
	fmt.Printf("---- srcMac: [%s], os: [%s] ----\n", string(srcMac), hostOS)
	ethOS[string(srcMac)] = hostOS

	// write log to file
	str := fmt.Sprintf("srcMac: [%s], os: [%s]\n", string(srcMac), hostOS)
	fd.WriteString(str)
	hostOS = ""
}
