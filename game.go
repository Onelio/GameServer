package main

import (
	"strings"
	"fmt"
	"strconv"
)

func printInsideRoom(client *Client) {
	client.SendPacket(fmt.Sprintf("*Estas en sala %s", client.room))

}

func printWaitToPlay(room *Room) {
	room.SendPacket("")
	room.SendPacket("[Sala] Contrincantes localizados, esperando a que esten listos...")
	room.users[0].SendPacket("[Mensaje] Di ready cuando estes listo")
	room.users[1].SendPacket("[Mensaje] Di ready cuando estes listo")
}

func printGame(room *Room) {
	l1, l2, l3 := getTable(room)
	room.SendPacket(l1)
	room.SendPacket(l2)
	room.SendPacket(l3)
}

func printWinner(room *Room, i int) {
	room.playing = false
	client := room.users[i]
	client.vcount++
	room.SendPacket("")
	room.SendPacket(fmt.Sprintf("[Estado] Victoria de %s!!", client.name))
	room.SendPacket(fmt.Sprintf("[Recuento de Victorias] %s: %s || %s: %s",
		room.users[0].name, strconv.Itoa(room.users[0].vcount), room.users[1].name, strconv.Itoa(room.users[1].vcount)))

	if i != 0 {
		room.SendPacket("[Estado] Inviertiendo posiciones!")
		c0 := room.users[0]
		room.users[0] = room.users[1]
		room.users[1] = c0
	}

	CheckForNextGame(room)
}

func CheckForNextGame(room *Room) {
	var winner, looser *Client
	if room.users[0].vcount == 3 {
		winner = room.users[0]
		looser = room.users[1]
	} else if room.users[1].vcount == 3 {
		winner = room.users[1]
		looser = room.users[0]
	} else {
		room.GameStart()
		return
	}

	room.SendPacket(fmt.Sprintf("[Estado] Enhorabuena %s has ganado 3 veces a %s",
		winner.name, looser.name))
	room.SendPacket("[Estado] Concluyendo la partida general!")
	looser.vcount = 0
	winner.vcount = 0

	if len(room.users) > 2 {
		room.SendPacket(fmt.Sprintf("[Sala] Retirando a %s del juego y añadiendo a %s",
			looser.name, room.users[2].name))
		room.QuitClient(looser)
		room.users = append(room.users, looser)
	}
	printWaitToPlay(room)
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
	if winner == -1 { //No winners
		room.SendPacket("[Estado] Partida terminada en empate!")
		room.SendPacket("[Sala] Empezando otra vez...")
		room.GameStart()
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
	room.SendPacket(fmt.Sprintf("[Estado] Turno de %s", client.name))
	client.SendPacket(" Como deseas seguir? Escribe play <numCasilla> para continuar")
	client.SendPacket("  Recuerda que se cuenta desde 1-9 de izquierda a derecha")
}

func (room *Room) GameStart() {
	room.playing = true
	room.turn = 0
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
		client.SendPacket("[Estado] Esperando jugadores para empezar la partida...")
		if len(room.users) > 1 {
			printWaitToPlay(room)
		}
	} else {
		client.SendPacket("[Sala] Partida en curso...")
	}
}

func (client *Client) LeaveRoom() {
	if client.room != "" {
		room := rooms[client.room]

		result := room.QuitClient(client)
		if result > -1 {
			client.room = ""
			room.SendPacket(fmt.Sprintf("El usuario %s ha abandonado la sala", client.name))

			//If user is player and is actually playing
			if room.playing && result < 2 {
				room.playing = false
				room.SendPacket("[Estado] Juego Cancelado!")
				room.CheckGameState(room.users[0])
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
	case "ready":
		client.SetReady()
	case "play":
		room := rooms[client.room]
		if room.users[room.turn] == client {
			room.PlayTurn(s[1])
		}
	default:
		printInsideRoom(client)
	}
}
