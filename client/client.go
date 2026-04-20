package client

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"unicode/utf8"

	"adbr.xx/anomchat/encryption"
	"adbr.xx/anomchat/terminal"
	"adbr.xx/anomchat/utils"
	"golang.org/x/term"
)

var inputSaved string
var state *term.State

// ----------------------
// CIFRADO + FRAMING
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
// INPUT DEL USUARIO
// ----------------------

func HandleUserInputs(conn net.Conn, chat utils.ChatInfo) {
	buf := make([]byte, 1)

	fmt.Print("> ")
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			terminal.RestoreTerminal(state)
			panic(err)
		}
		if n > 0 {
			b := buf[0]

			if b == 13 { // ENTER
				if len(inputSaved) > 0 {
					fmt.Println("\r")
					fmt.Print("> ")
					SendEncrypted(conn, []byte(inputSaved), []byte(chat.Key))
					inputSaved = ""
				}
				continue
			} else if b == 127 || b == 8 { // BACKSPACE
				if len(inputSaved) > 0 {
					_, size := utf8.DecodeLastRuneInString(inputSaved)
					inputSaved = inputSaved[:len(inputSaved)-size]
					fmt.Print("\b \b")
				}
			} else if b >= 32 && b <= 126 {
				inputSaved += string(b)
				fmt.Print(string(b))
			} else if b == '\x03' { // CTRL+C
				terminal.RestoreTerminal(state)
				os.Exit(0)
			}
		}
	}
}

// ----------------------
// CLIENTE TCP
// ----------------------

func InitializeClient(chat utils.ChatInfo) {
	s, err := net.Dial("tcp", chat.Host)
	if err != nil {
		terminal.RestoreTerminal(state)
		panic(err)
	}

	// Enviar clave cifrada para autenticación
	SendEncrypted(s, []byte(chat.Key), []byte(chat.Key))

	state = terminal.SetupTerminal()
	if state == nil {
		panic("Failed to set up terminal")
	}
	defer func() { s.Close(); terminal.RestoreTerminal(state) }()

	// LECTOR DE MENSAJES
	go func() {
		for {
			msgBytes, err := ReadEncrypted(s, []byte(chat.Key))
			if err != nil {
				terminal.RestoreTerminal(state)
				panic(err)
			}

			message := string(msgBytes)

			if message == "KICKED" {
				terminal.RestoreTerminal(state)
				os.Exit(0)
			}

			fmt.Print("\r\033[K")
			fmt.Println(message, "\r")
			fmt.Print("> ", inputSaved)
		}
	}()

	utils.ShowChatInfo(chat, false)

	go HandleUserInputs(s, chat)
	select {}
}
