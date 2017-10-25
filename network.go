package main

import (
	"net"
	"strings"
	"bufio"
)

type Client struct {
	conn	net.Conn
	ochan	chan string
	name	string
	room	string
	ready	bool
	vcount	int
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
		if err != nil {
			break
		}

		line = strings.TrimSuffix(line, "\n")
		if line != "" {
			server.Handler(client, line)
		}
	}
	RemoveClient(client)
}

func (client *Client) outputChannel(writer *bufio.Writer) {
	for data := range client.ochan {
		writer.WriteString(data)
		writer.Flush()
	}
}

func (client *Client) SendPacket(str string) {
	if !strings.Contains(str, "\n") {
		str += "\n"
	}
	client.ochan <- str
}

func (room *Room) SendPacket(string string) {
	for _, client := range room.users {
		client.SendPacket(string)
	}
}

func NewClient(conn net.Conn) (*Client) {
	client := &Client{
		conn: conn,
		ochan: make(chan string),
	}

	go client.inputChannel(bufio.NewReader(client.conn))
	go client.outputChannel(bufio.NewWriter(client.conn))

	return client
}

func RemoveClient(client *Client) {
	client.conn.Close()
	//Remove from room
	client.LeaveRoom()

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