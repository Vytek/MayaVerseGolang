package main

import (
	"net"
	"strings"

	log "github.com/Masterminds/log-go"
	"github.com/lrita/cmap"
	"github.com/obsilp/rmnp"
	"github.com/rs/xid"
	"github.com/vmihailenco/msgpack/v5"
	"gitlab.com/rwxrob/uniq"
)

var guid xid.ID
var n cmap.Map[string, *rmnp.Connection]

type Messages struct {
	OpCode  byte
	Message string
}

func main() {
	//Unique ID
	guid = xid.New()
	log.Infof("Server uniqueID: %s\n", guid.String())

	server := rmnp.NewServer(":10001") //TODO: Add ini config for port and others

	server.ClientConnect = clientConnect
	server.ClientDisconnect = clientDisconnect
	server.ClientTimeout = clientTimeout
	server.ClientValidation = validateClient
	server.PacketHandler = handleServerPacket

	server.Start()
	log.Infof("Server started")

	select {}
}

func clientConnect(conn *rmnp.Connection, data []byte) {
	log.Infof("Client connection with:", data)

	UniqueID := uniq.Hex(18)
	if data[0] != 0 {
		conn.Disconnect([]byte("not allowed"))
	} else {
		//Add new client connected
		n.Store(UniqueID+":"+conn.Addr.String(), conn)
		conn.SendReliableOrdered([]byte(UniqueID + ":" + conn.Addr.String()))
	}
}

func clientDisconnect(conn *rmnp.Connection, data []byte) {
	log.Infof("client disconnect with:", data)
	//Parse Message received

	//Delete the client connected from cmap
}

func clientTimeout(conn *rmnp.Connection, data []byte) {
	log.Infof("Client timeout with:", data)
	//Delete the client Timeouted
}

func validateClient(addr *net.UDPAddr, data []byte) bool {
	//Parse Message received
	var MessageReceived Messages
	err := msgpack.Unmarshal(data, &MessageReceived)
	if err != nil {
		panic(err)
	}
	log.Infof(MessageReceived.Message)
	s := strings.Split(string(MessageReceived.Message), ":")
	log.Infof(s[0])
	return len(data) == 3
}

func handleServerPacket(conn *rmnp.Connection, data []byte, channel rmnp.Channel) {
	str := string(data)
	log.Infof("'"+str+"'", "from", conn.Addr.String(), "on channel", channel)

	if str == "ping" {
		conn.SendReliableOrdered([]byte("pong"))
		conn.Disconnect([]byte("session end"))
	}
}
