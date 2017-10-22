package main

import (
	"net"
	"strings"
	"bufio"
	"time"
)

type Client struct {
	conn	net.Conn
	ochan	chan string
	name	string
	room	string
	gamer	bool
}

type Server struct {
	Host 	string
	Port	string
	Handler func(*Client, string)
	Clients	[]*Client
}

var server *Server

func (client *Client) inputChannel(reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')
		if err == nil {
			line := strings.TrimSuffix(line, "\n")
			if line != "" {
				server.Handler(client, line)
			}
		} else {
			break
		}
	}
	RemoveClient(client)
}

func (client *Client) outputChannel(writer *bufio.Writer) {
	for data := range client.ochan {
		time.Sleep(1) //Prevent bug
		writer.WriteString(data)
		writer.Flush()
	}
}

func (client *Client) SendPacket(string string) {
	if !strings.Contains(string, "\n") {
		string += "\n"
	}
	client.ochan <- string
}

func (room *Room) SendPacket(string string) {
	if !strings.Contains(string, "\n") {
		string += "\n"
	}

	for _, client := range room.users {
		client.SendPacket(string)
	}
}

func NewClient(conn net.Conn) (*Client) {
	client := &Client{ conn: conn, ochan: make(chan string)}
	go client.inputChannel(bufio.NewReader(client.conn))
	go client.outputChannel(bufio.NewWriter(client.conn))
	return client
}

func RemoveClient(client *Client) {
	client.conn.Close()
	//Remove from room
	if client.room != "" {
		room := rooms[client.room]
		for i, eclient := range room.users {
			if eclient == client {
				room.users = append(room.users[:i], room.users[i+1:]...)
				room.SendPacket("El usuario " + client.name + " ha abandonado la sala")
			}
		}
		if len(room.users) < 1 {
			delete(rooms, room.name)
		}
	}

	//Remove from client list
	for i, eclient := range server.Clients {
		if eclient == client {
			server.Clients = append(server.Clients[:i], server.Clients[i+1:]...)
			client = nil
		}
	}
}

func (server *Server) StartListening() (err error) {
	listener, _ := net.Listen("tcp", server.Host + ":" + server.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			//If cannot be accepted we get rid of it
			continue
		}
		client := NewClient(conn)
		server.Clients = append(server.Clients, client)

		//Welcome
		client.SendPacket("Bienvenido al 3 en raya online!")
		client.SendPacket("Por favor introduce tu nombre a continuacion:")
	}
}