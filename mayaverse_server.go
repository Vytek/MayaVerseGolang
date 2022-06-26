package main

import (
	"fmt"
	"net"

	"github.com/lrita/cmap"
	"github.com/obsilp/rmnp"
	"github.com/rs/xid"
	"gitlab.com/rwxrob/uniq"
)

var guid xid.ID
var n cmap.Map[string, *rmnp.Connection]

type Messages struct {
	OpCode  int
	Message string
}

func main() {
	//Unique ID
	guid = xid.New()
	fmt.Printf("%s\n", guid.String())

	server := rmnp.NewServer(":10001")

	server.ClientConnect = clientConnect
	server.ClientDisconnect = clientDisconnect
	server.ClientTimeout = clientTimeout
	server.ClientValidation = validateClient
	server.PacketHandler = handleServerPacket

	server.Start()
	fmt.Println("server started")

	select {}
}

func clientConnect(conn *rmnp.Connection, data []byte) {
	fmt.Println("client connection with:", data)

	if data[0] != 0 {
		conn.Disconnect([]byte("not allowed"))
	} else {
		//Add new client connected
		n.Store(uniq.Hex(18)+":"+conn.Addr.String(), conn)
	}
}

func clientDisconnect(conn *rmnp.Connection, data []byte) {
	fmt.Println("client disconnect with:", data)
	//Delete the client connected
}

func clientTimeout(conn *rmnp.Connection, data []byte) {
	fmt.Println("client timeout with:", data)
	//Delete the client Timeouted
}

func validateClient(addr *net.UDPAddr, data []byte) bool {
	return len(data) == 3
}

func handleServerPacket(conn *rmnp.Connection, data []byte, channel rmnp.Channel) {
	str := string(data)
	fmt.Println("'"+str+"'", "from", conn.Addr.String(), "on channel", channel)

	if str == "ping" {
		conn.SendReliableOrdered([]byte("pong"))
		conn.Disconnect([]byte("session end"))
	}
}
