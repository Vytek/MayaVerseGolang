package main

import (
	"fmt"

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

	select {}
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
