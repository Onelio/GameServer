package main

import (
	"strings"
	"fmt"
)

func PacketReceived(client *Client, packet string) {
	packet = strings.Replace(packet, "\r", "", -1)
	if client.name == "" {
		client.name = packet
		client.SendPacket(fmt.Sprintf("Conectado como %s", client.name))
	}

	if client.room == "" {
		handleChannelPacket(client, packet)
	} else {
		handleRoomPacket(client, packet)
	}
}

//TODO HANDLE REMOTE SERVER CONTROLS

func main() {
	println("Starting game server")

	rooms = make(map[string]*Room)
	server = &Server{
		Host:    "localhost",
		Port:    "8081",
		Handler: PacketReceived,
	}

	println(fmt.Sprintf("Listening to %s:%s ...", server.Host, server.Port))
	server.StartListening()
}
