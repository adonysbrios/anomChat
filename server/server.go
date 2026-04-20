package server

import (
	"encoding/binary"
	"io"
	"net"
	"sync"

	"adbr.xx/anomchat/encryption"
	"adbr.xx/anomchat/utils"
)

// ----------------------
// CIFRADO y FRAMING
// ----------------------

func SendEncrypted(conn net.Conn, plaintext []byte, key []byte) error {
	ciphertext, err := encryption.EncryptMessage(plaintext, key)
	if err != nil {
		return err
	}

	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(ciphertext)))

	if _, err := conn.Write(lenBuf); err != nil {
		return err
	}
	if _, err := conn.Write(ciphertext); err != nil {
		return err
	}

	return nil
}

func ReadEncrypted(conn net.Conn, key []byte) ([]byte, error) {
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return nil, err
	}

	msgLen := binary.BigEndian.Uint32(lenBuf)
	msgBuf := make([]byte, msgLen)

	if _, err := io.ReadFull(conn, msgBuf); err != nil {
		return nil, err
	}

	return encryption.DecryptMessage(msgBuf, key)
}

// ----------------------
// Chat server
// ----------------------

var Chat utils.ChatInfo

type Message struct {
	owner   net.Conn
	message string
}

var broadcastChannel = make(chan Message, 100)
var clientsMutex = &sync.Mutex{}
var clients = make(map[net.Conn]string)

func handleConnection(conn net.Conn) {
	// Recibir clave cifrada
	keyBytes, err := ReadEncrypted(conn, []byte(Chat.Key))
	if err != nil {
		conn.Close()
		return
	}

	key := string(keyBytes)
	if key != Chat.Key {
		SendEncrypted(conn, []byte("KICKED"), []byte(Chat.Key))
		conn.Close()
		return
	}

	clientsMutex.Lock()
	clients[conn] = utils.GenerateRandomUsername()
	clientsMutex.Unlock()

	// LECTOR DE MENSAJES
	go func() {
		defer func() {
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			conn.Close()
		}()

		for {
			msgBytes, err := ReadEncrypted(conn, []byte(Chat.Key))
			if err != nil {
				return
			}

			message := string(msgBytes)
			broadcastChannel <- Message{owner: conn, message: message}
		}
	}()
}

func BroadcastMessages() {
	for message := range broadcastChannel {
		clientsMutex.Lock()
		for conn, username := range clients {
			if conn != message.owner {
				SendEncrypted(conn, []byte(username+": "+message.message), []byte(Chat.Key))
			}
		}
		clientsMutex.Unlock()
	}
}

func InitializeTCPServer(port, key string) {
	s, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	go BroadcastMessages()

	Chat = utils.ChatInfo{Host: s.Addr().String(), Key: key}
	utils.ShowChatInfo(Chat, true)

	defer s.Close()

	for {
		conn, err := s.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
