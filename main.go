package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	// Set the network interface and filter for capturing packets
	iface := "\\Device\\NPF_Loopback"
	filter := "tcp and (port 54321 or port 12345)"
	clientSrcPort := ""
	// Open the network interface for packet capture
	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal("Error opening interface:", err)
	}
	defer handle.Close()
	// Set a BPF filter to capture only the desired packets
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal("Error setting BPF filter:", err)
	}

	// Create a packet source from the handle
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Loop to process incoming packets
	for packet := range packetSource.Packets() {
		applicationLayer := packet.ApplicationLayer()

		if applicationLayer != nil {
			//fmt.Printf("Application layer: %s\n", applicationLayer.Payload())

			if strings.Contains(string(applicationLayer.Payload()), "GET") {

				fmt.Printf("Request: %s\n", applicationLayer.Payload())
				sendRequest(string(applicationLayer.Payload()))
				tcpLayer := packet.Layer(layers.LayerTypeTCP)
				if tcpLayer != nil {
					tcp, _ := tcpLayer.(*layers.TCP)

					clientSrcPort = tcp.SrcPort.String()
					//fmt.Println(clientSrcPort)
				}

			} else if strings.Contains(string(applicationLayer.Payload()), "Response") {

				fmt.Printf("Response: %s\n", applicationLayer.Payload())

				// Send the modified packet forward
				tcpLayer := packet.Layer(layers.LayerTypeTCP)
				if tcpLayer != nil {
					tcp, _ := tcpLayer.(*layers.TCP)

					// Check if the UDP packet is destined for the specific port

					forwardResponse(string(applicationLayer.Payload()), tcp.SrcPort.String(), clientSrcPort)
				}
			}
		}
	}
}

func sendRequest(data string) {
	// Set the local address and port to which you want to send the string
	dstAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatal("Error resolving address:", err)
	}

	// Create a TCP packet
	TCPLayer := &layers.TCP{
		SrcPort: layers.TCPPort(0), // Let the OS choose the source port
		DstPort: layers.TCPPort(dstAddr.Port),
	}

	// Set the payload (your string message)
	//payload := []byte(data)

	// Create a buffer to store the entire packet
	buffer := gopacket.NewSerializeBuffer()

	// Serialize the layers and payload into the buffer
	err = gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{},
		TCPLayer,
	)
	if err != nil {
		log.Fatal("Error serializing packet:", err)
	}

	// Get the raw bytes from the buffer
	packetData := buffer.Bytes()

	// Send the packet to the destination address
	conn, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		log.Fatal("Error creating connection:", err)
	}
	defer conn.Close()

	_, err = conn.Write(packetData)
	if err != nil {
		log.Fatal("Error sending packet:", err)
	}

	fmt.Println("Packet sent successfully.")
	time.Sleep(time.Second) // Give the packet some time to be sent before the program exits
}

func forwardResponse(data string, srcPort string, DstPort string) {
	// Set the local address and port to which you want to send the string
	dstAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+DstPort)
	if err != nil {
		log.Fatal("Error resolving address:", err)
	}
	srcPort = strings.Split(srcPort, "(")[0]
	srcPortNumber, err := strconv.Atoi(srcPort)

	if err != nil {
		log.Fatal("Error serializing packet:", err)
	}

	// Create a TCP packet
	TCPLayer := &layers.TCP{
		SrcPort: layers.TCPPort(srcPortNumber), // Let the OS choose the source port
		DstPort: layers.TCPPort(dstAddr.Port),
	}

	// Set the payload (your string message)
	payload := []byte(data)

	// Create a buffer to store the entire packet
	buffer := gopacket.NewSerializeBuffer()

	// Serialize the layers and payload into the buffer
	err = gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{},
		TCPLayer,
		gopacket.Payload(payload),
	)
	if err != nil {
		log.Fatal("Error serializing packet:", err)
	}

	// Get the raw bytes from the buffer
	packetData := buffer.Bytes()

	// Send the packet to the destination address
	conn, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		log.Fatal("Error creating connection:", err)
	}
	defer conn.Close()

	_, err = conn.Write(packetData)
	if err != nil {
		log.Fatal("Error sending packet:", err)
	}

	fmt.Println("Packet Forwarded successfully.")
	time.Sleep(time.Second) // Give the packet some time to be sent before the program exits
}

func sendPacket(packet []byte, dstPort string) error {
	// Create a TCP connection to send the packet
	conn, err := net.Dial("tcp", "127.0.0.1:"+dstPort) // Change this to the destination address
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the packet
	_, err = conn.Write(packet)
	return err
}
