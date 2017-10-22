package main

import "strings"

func PacketReceived(client *Client, packet string) {
	packet = strings.TrimSuffix(packet, "\r")
	if client.name == "" {
		client.name = packet
		client.SendPacket("Conectado como " + client.name)
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
	server = &Server{ Host: "localhost", Port: "8081", Handler: PacketReceived }
	println("Listening to " + server.Host + ":" + server.Port + " ...")
	server.StartListening()
}
