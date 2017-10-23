package main

import (
	"strings"
	"fmt"
	"strconv"
)

func printInsideRoom(client *Client) {
	mode := "Jugador"
	if !client.gamer {
		mode = "Espectador"
	}
	client.SendPacket(fmt.Sprintf("*Sala %s Mode: %s", client.room, mode))

}

func printGame(room *Room) {
	l1, l2, l3 := getTable(room)
	room.SendPacket(l1)
	room.SendPacket(l2)
	room.SendPacket(l3)
}

func printWinner(room *Room, i int) {
	client := room.users[i]
	client.vcount++
	room.SendPacket("")
	room.SendPacket("[Estado] Victoria de " + client.name + "!!")
	room.SendPacket("[Recuento de Victorias]")
	room.SendPacket(" - " + room.users[0].name + ": " + strconv.Itoa(room.users[0].vcount))
	room.SendPacket(" - " + room.users[1].name + ": " + strconv.Itoa(room.users[1].vcount))

	if i != 0 {
		room.SendPacket("[Estado] Inviertiendo posiciones!")
		c0 := room.users[0]
		room.users[0] = room.users[1]
		room.users[1] = c0
	}
	room.GameStart()
}

func (room *Room) PlayTurn(str string) {

	loc, err := strconv.Atoi(str)
	if str == "" || err != nil || loc > 9 {
		return
	}

	if room.cells[loc-1] != 0 { //Do not allow cheats
		return
	}
	room.cells[loc-1] = room.turn + 1 //To not count 0
	printGame(room)

	winner := findWinner(room)
	if winner > 0 {
		printWinner(room, winner-1) //Return to count from 0
		return
	}
	room.turn++
	if room.turn > 1 {
		room.turn = 0
	}
	room.GiveTurn()
}

func (room *Room) GiveTurn() {
	client := room.users[room.turn]
	room.SendPacket("[Estado] Turno de " + client.name)
	client.SendPacket(" Como deseas seguir? Escribe play <numCasilla> para continuar")
	client.SendPacket("  Recuerda que se cuenta desde 1-9 de izquierda a derecha")
}

func (room *Room) GameStart() {
	room.playing = true
	room.cells = make([]int, 9)

	room.SendPacket("")
	room.SendPacket("[Estado] ¡Que comienze la partida!")
	room.SendPacket("[Contrincantes]")
	room.SendPacket(" - Jugando como [X] tenemos a " + room.users[0].name)
	room.SendPacket(" - Jugando como [0] tenemos a " + room.users[1].name)
	printGame(room)
	room.SendPacket("")

	room.GiveTurn()
}

func (room *Room) CheckGameState(client *Client) {
	if !room.playing {
		if len(room.users) < 2 {
			client.SendPacket("[Estado] Esperando jugadores para empezar la partida...")
		} else {
			room.users[0].vcount = 0
			room.users[1].vcount = 0
			room.GameStart()
		}
	}
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
	case "play":
		room := rooms[client.room]
		if room.users[room.turn] == client {
			room.PlayTurn(s[1])
		}
	default:
		printInsideRoom(client)
	}
}
