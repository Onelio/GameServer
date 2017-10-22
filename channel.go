package main

import (
	"strconv"
	"strings"
)

type Room struct {
	name 	string
	mode	string
	users	[]*Client
	playing	bool
}

var rooms map[string]*Room

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
		client.SendPacket("[" + strconv.Itoa(count) + "]" + room.name + " Users:" + strconv.Itoa(users) + " Mode:" + room.mode)
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
	if len(room.users) < 3 {
		client.gamer = true
	} else {
		client.gamer = false
	}
	room.users = append(room.users, client)

	printInsideRoom(client)
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
		break
	case "enter":
		if len(s) < 2 || s[1] == "" { //Blank not valid
			client.SendPacket("No especificado")
			return
		}
		enterRoom(client, s[1])
		break
	default:
		printChannelData(client)
	}
}
