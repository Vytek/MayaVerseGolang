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
	log.Infof("Client connection with: %s\n", data)

	UniqueID := uniq.Hex(18)
	if data[0] != 0 {
		conn.Disconnect([]byte("not allowed"))
	} else {
		//Add new client connected
		n.Store(UniqueID+":"+conn.Addr.String(), conn)
		b, err := msgpack.Marshal(&Messages{OpCode: 1, Message: "cid:" + UniqueID + ":" + conn.Addr.String()})
		if err != nil {
			log.Errorf("Can't create MessagePack OpCode 1 Message")
		} else {
			conn.SendReliableOrdered(b)
		}
	}
}

func clientDisconnect(conn *rmnp.Connection, data []byte) {
	log.Infof("Client disconnect with: %s\n", data)
	//Parse Message received
	var MessageReceived Messages
	err := msgpack.Unmarshal(data, &MessageReceived)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		return
	}
	log.Infof(MessageReceived.Message)
	s := strings.Split(string(MessageReceived.Message), ":")
	if MessageReceived.OpCode == 2 {
		//Delete the client connected from cmap
		if s[0] == "cld" {
			n.Delete(s[1])
		} else {
			log.Errorf("Not cld in Message command")
		}
	} else {
		log.Errorf("Not Opcode 2 in Message")
	}
}

func clientTimeout(conn *rmnp.Connection, data []byte) {
	log.Infof("Client timeout with: %s\n", data)
	//Delete the client Timeouted
	var ClientToDelete string
	n.Range(func(key string, value *rmnp.Connection) bool {
		k, v := key, value
		if v.Addr.String() == conn.Addr.String() {
			ClientToDelete = k
		}
		return true
	})
	n.Delete(ClientToDelete)
}

func validateClient(addr *net.UDPAddr, data []byte) bool {
	//Parse Message received
	var MessageReceived Messages
	err := msgpack.Unmarshal(data, &MessageReceived)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		return false
	}
	log.Infof(MessageReceived.Message)
	s := strings.Split(string(MessageReceived.Message), ":")
	if MessageReceived.OpCode == 0 {
		if s[0] == "lng" {
			return true
			//Check login and password using scrypt
		} else {
			log.Errorf("Not lng in Message command")
			return false
		}
	} else {
		log.Errorf("Not Opcode 1 in Message")
		return false
	}
}

func handleServerPacket(conn *rmnp.Connection, data []byte, channel rmnp.Channel) {
	str := string(data)
	log.Infof("'"+str+"'", "from", conn.Addr.String(), "on channel", channel)

	if str == "ping" {
		conn.SendReliableOrdered([]byte("pong"))
		conn.Disconnect([]byte("session end"))
		//Delete client disconnected from cmap
		var ClientToDelete string
		n.Range(func(key string, value *rmnp.Connection) bool {
			k, v := key, value
			if v.Addr.String() == conn.Addr.String() {
				ClientToDelete = k
			}
			return true
		})
		n.Delete(ClientToDelete)
	}
}
