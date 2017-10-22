package main

import "strings"

func printInsideRoom(client *Client) {
	mode := "Jugador"
	if !client.gamer {
		mode = "Espectador"
	}
	client.SendPacket("*Sala " + client.room + " Mode: " + mode)
}

func handleRoomPacket(client *Client, packet string) {
	s := strings.Split(packet, " ")

	switch s[0] {
	case "leave":
		break
	case "enter":
		break
	default:
		printInsideRoom(client)
	}
}
