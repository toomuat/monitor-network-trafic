package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	windowsCounter uint64
	androidCounter uint64
	iosCounter     uint64

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
)

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
