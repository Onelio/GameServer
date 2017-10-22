package main

import (
	"strings"
	"fmt"
)

func printInsideRoom(client *Client) {
	mode := "Jugador"
	if !client.gamer {
		mode = "Espectador"
	}
	client.SendPacket(fmt.Sprintf("*Sala %s Mode: %s", client.room, mode))
}

func (client *Client) LeaveRoom() {
	if client.room != "" {
		room := rooms[client.room]
		for i, eclient := range room.users {
			if eclient == client {
				room.users = append(room.users[:i], room.users[i+1:]...)
				client.room = ""
				room.SendPacket("El usuario " + client.name + " ha abandonado la sala")
			}
		}
		if len(room.users) < 1 {
			delete(rooms, room.name)
		}
	}
}

func handleRoomPacket(client *Client, packet string) {
	s := strings.Split(packet, " ")

	switch s[0] {
	case "leave":
		client.LeaveRoom()
		printChannelData(client)
	case "say":
		packet = strings.Replace(packet, "say ", "", -1)
		if packet != "" {
			room := rooms[client.room]
			room.SendPacket(client.name + ": " + packet)
		}
	default:
		printInsideRoom(client)
	}
}
