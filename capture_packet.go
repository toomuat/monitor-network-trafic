package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	windowsCounter uint64
	androidCounter uint64
	iosCounter     uint64

	messagesMap = map[string]*Message{
		"Windows": &Message{"Windows", 0},
		"Android": &Message{"Android", 0},
		"iOS":     &Message{"iOS", 0},
	}
	messages  []Message
	messages2 = []*Message{
		&Message{"Windows", 0},
		&Message{"Android", 0},
		&Message{"iOS", 0},
	}

	// ethOS map[net.HardwareAddr]string
	ethOS map[string]string

	// variable for packet capture
	eth         string = "eth1"
	snapshotLen int32  = 1024
	promiscuous bool   = true
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

	// matchWindows = regexp.MustCompile(`MSFT`)
	// matchAndroid = regext.MustCompile(`android`)
	// matchiOS = regext.MustCompile(`iPhone`)
)

func sendCounter(w http.ResponseWriter, r *http.Request) {
	upgrade, _ := upgrader.Upgrade(w, r, nil)

	// this could be deleted
	// go func() {
	// for {
	// count packet of each OS
	// time.Sleep(1 * time.Microsecond)
	// }
	// 	capturePacket()
	// }()

	for {
		// send number of packet every 1 second
		time.Sleep(1 * time.Second)

		messages = messages[:0] // clear slice
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

func packetCounter() {
	for {
		// count packet of each OS
		time.Sleep(1 * time.Microsecond)
	}
}

func capturePacket() {
	handle, err = pcap.OpenLive(eth, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packets := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packets.Packets() {
		// fmt.Println(packet)
		countPacket(packet)
	}
}

func countPacket(packet gopacket.Packet) {
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
		dhcpFingerprinting(srcMac, dhcpPacket)
		fmt.Println(dhcpLayer)
	}
}

func dhcpFingerprinting(srcMac string, dhcpPacket *layers.DHCPv4) {
	fmt.Println("Operation: ", dhcpPacket.Operation)
	fmt.Println("Options: ")
	os_name = os_name[:0]
	for _, option := range dhcpPacket.Options {
		fmt.Println(option.String())
		if option.Type == layers.DHCPOptHostname {
			os_name = append(os_name, string(option.Data))
		} else if option.Type == layers.DHCPOptClassID {
			os_name = append(os_name, string(option.Data))
		}
	}

	for _, os := range os_name {
		switch {
		case strings.Contains(os, "MSFT"):
			host_os = "Windows"
		case strings.Contains(os, "android"):
			host_os = "Android"
		case strings.Contains(os, "iPhone"):
			host_os = "iOS"
		default:
			fmt.Println("nothing matched")
		}
	}

	ethOS[host_os] = srcMac
}
