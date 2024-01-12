package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var DeviceName = "\\Device\\NPF_Loopback"
var DeviceFound = false

func main() {
	// find all network interfaces:
	devices, err := pcap.FindAllDevs()

	if err != nil {
		log.Panicln("Unable to fetch network interfaces")
	}

	// if device is found, continue

	for _, device := range devices {
		//fmt.Println(device.Name)
		if device.Name == DeviceName {
			DeviceFound = true
		}
	}

	if !DeviceFound {
		log.Panicln("Desired Device not found")
	}

	//Open live capture on network interface:
	handle, err := pcap.OpenLive(DeviceName, 1600, false, pcap.BlockForever)

	if err != nil {
		fmt.Println(err)
		log.Panicln("Unable to open handle on the device")
	}

	defer handle.Close()

	// Apply filters:

	if err := handle.SetBPFFilter("tcp and port 12345"); err != nil {
		fmt.Println(err)
		log.Panicln("Unable to filter packets")
	}

	//Display the filtered packts:

	source := gopacket.NewPacketSource(handle, handle.LinkType())

	packets := source.Packets()
	for packet := range packets {
		//fmt.Println(packet)

		applicationLayer := packet.ApplicationLayer()
		if applicationLayer != nil {
			fmt.Println("Application layer/Payload found.")
			fmt.Printf("Message from client: %s\n", applicationLayer.Payload())

		}
	}
}
