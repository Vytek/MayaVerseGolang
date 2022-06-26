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

	b, err := msgpack.Marshal(&Messages{OpCode: 0, Message: "lng:" + ""})
	if err != nil {
		panic(err)
	}
	client.ConnectWithData(b)

	select {}
}

func serverConnect(conn *rmnp.Connection, data []byte) {
	fmt.Println("connected to server")
	conn.SendReliableOrdered([]byte("ping"))
	ServerConnection = conn
}

func serverDisconnect(conn *rmnp.Connection, data []byte) {
	fmt.Println("disconnected from server:", string(data))
}

func serverTimeout(conn *rmnp.Connection, data []byte) {
	fmt.Println("server timeout")
}

func handleClientPacket(conn *rmnp.Connection, data []byte, channel rmnp.Channel) {
	fmt.Println("'"+string(data)+"'", "on channel", channel)
}

func SendMessage(Message string) {
	ServerConnection.SendOnChannel(1, []byte(Message))
}
