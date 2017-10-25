package main

import (
	"strconv"
	"strings"
	"fmt"
)

type Room struct {
	name 	string
	mode	string
	users	[]*Client
	playing	bool
	turn	int
	cells	[]int
}

var rooms map[string]*Room

func (room *Room) QuitClient(client *Client) (int) {
	for i, eclient := range room.users {
		if eclient == client {
			room.users = append(room.users[:i], room.users[i+1:]...)
			return i
		}
	}
	return -1
}

func (client *Client) SetReady() {
	room := rooms[client.room]
	client.ready = true

	if room.users[0].ready && room.users[1].ready {
		room.users[0].ready = false
		room.users[1].ready = false
		room.SendPacket("[Sala] Todos listos!")
		room.users[0].vcount = 0
		room.users[1].vcount = 0
		room.GameStart()
	} else {
		room.SendPacket(fmt.Sprintf("[Sala] Contrincante %s esta listo!", client.name))
	}
}

func printChannelData(client *Client) {
	client.SendPacket("[Panel de Seleccion de Sala]")
	client.SendPacket(" - Para crear una sala has de enviar create <nombre_sin_espacios>")
	client.SendPacket(" - Para entrar a una sala has de enviar enter <nombre_sin_espacios>")
	client.SendPacket("")
	client.SendPacket("[Lista de canales existentes]")

	//Start printing
	count := 0
	for _, room := range rooms {
		count++
		users := len(room.users)
		client.SendPacket(fmt.Sprintf("[%s] %s Users: %s Mode: %s",
			strconv.Itoa(count), room.name, strconv.Itoa(users), room.mode))
	}
	client.SendPacket("")
}

func enterRoom(client *Client, name string) {
	room, exist := rooms[name]
	if !exist {
		client.SendPacket("Error, a la sala que deseas entrar no existe!")
		return
	}
	client.room = name
	client.ready = false
	room.users = append(room.users, client)

	room.SendPacket(fmt.Sprintf("El usuario %s ha entrado a la sala", client.name))
	room.CheckGameState(client)
}

func createRoom(client *Client, name string) {
	_, exist := rooms[name]
	if exist {
		client.SendPacket("Error, la sala que deseas crear ya existe!")
		return
	}
	room := &Room{
		name: name,
		mode: "Best of 3",
	}
	rooms[name] = room

	enterRoom(client, name)
}

func handleChannelPacket(client *Client, packet string) {
	s := strings.Split(packet, " ")

	switch s[0] {
	case "create":
		if len(s) < 2 || s[1] == "" { //Blank not valid
			client.SendPacket("No especificado")
			return
		}
		createRoom(client, s[1])
	case "enter":
		if len(s) < 2 || s[1] == "" { //Blank not valid
			client.SendPacket("No especificado")
			return
		}
		enterRoom(client, s[1])
	default:
		printChannelData(client)
	}
}
