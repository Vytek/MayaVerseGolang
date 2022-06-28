package main

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell/v2"
	"github.com/obsilp/rmnp"
	"github.com/vmihailenco/msgpack/v5"
)

type Messages struct {
	OpCode  byte
	Message string
}

var ServerConnection *rmnp.Connection

func main() {
	client := rmnp.NewClient("127.0.0.1:10001")

	client.ServerConnect = serverConnect
	client.ServerDisconnect = serverDisconnect
	client.ServerTimeout = serverTimeout
	client.PacketHandler = handleClientPacket

	b, err := msgpack.Marshal(&Messages{OpCode: 0, Message: "lng:" + "login"})
	if err != nil {
		panic(err)
	} else {
		client.ConnectWithData(b)
	}

	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// display welcome info.
	shell.Println("MayaVerse Client Interactive Shell")

	// register a function for "greet" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "greet",
		Help: "greet user",
		Func: func(c *ishell.Context) {
			c.Println("Hello", strings.Join(c.Args, " "))
		},
	})

	// run shell
	shell.Run()
}

func serverConnect(conn *rmnp.Connection, data []byte) {
	fmt.Println("Connected to server with data: " + string(data))
	//Parse OpCode 1
	conn.SendReliableOrdered([]byte("ping"))
	ServerConnection = conn
}

func serverDisconnect(conn *rmnp.Connection, data []byte) {
	fmt.Println("Disconnected from server: " + string(data))
}

func serverTimeout(conn *rmnp.Connection, data []byte) {
	fmt.Println("Server timeout")
}

func handleClientPacket(conn *rmnp.Connection, data []byte, channel rmnp.Channel) {
	fmt.Println("'"+string(data)+"'", "on channel", channel)
}

func SendMessage(Message string) {
	ServerConnection.SendOnChannel(1, []byte(Message))
}
