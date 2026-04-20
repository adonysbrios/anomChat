<img width="882" height="155" alt="Captura desde 2026-04-19 14-08-12" src="https://github.com/user-attachments/assets/9b07a4b8-f69f-4157-b16d-16d15cfbd513" />


# anomChat

A secure, encrypted command-line chat application built in Go. anomChat enables anonymous peer-to-peer communication between a client and server using AES-256 encryption with automatic key derivation.

## Features

- **End-to-End Encryption**: All messages are encrypted using AES-256-GCM with SHA-256 key derivation
- **Anonymous Communication**: Lightweight and privacy-focused design
- **TCP-Based**: Reliable network communication
- **CLI Interface**: Simple command-line interface for both client and server modes
- **Auto-Generated Passwords**: Server generates secure random passwords for client authentication

## Project Structure

```
├── main.go              # Entry point for client/server mode selection
├── client/
│   └── client.go        # Client implementation with message handling
├── server/
│   └── server.go        # Server implementation with multi-client support
├── encryption/
│   └── encryption.go    # AES-256-GCM encryption and decryption logic
├── terminal/
│   └── terminal.go      # Terminal UI utilities
├── utils/
│   └── utils.go         # Helper functions and data structures
└── go.mod               # Go module definition
```

## Installation

### Prerequisites

- Go 1.26.2 or later

### Build

```bash
go build -o anomchat
```

## Usage

### Starting a Server

```bash
./anomchat s <port>
```

**Example:**
```bash
./anomchat s 8080
```

The server will:
- Start listening on the specified port
- Generate a random encryption key automatically
- Display the key for clients to use

### Connecting as a Client

```bash
./anomchat c <host> <key>
```

**Example:**
```bash
./anomchat c localhost:8080 your-encryption-key
```

Where:
- `<host>` is the server address (e.g., `localhost:8080` or `192.168.1.100:5000`)
- `<key>` is the encryption key provided by the server

## Security

- **Encryption**: AES-256-GCM (Galois/Counter Mode) ensures both confidentiality and authenticity
- **Key Derivation**: Server-provided keys are processed through SHA-256 to generate 256-bit cipher keys
- **Message Framing**: Length-prefixed messages prevent truncation attacks

## Dependencies

- `golang.org/x/term` - Terminal handling for cross-platform support
- `golang.org/x/sys` - System-level utilities

## Future Enhancements

- Commands

## License

MIT

## Contributing

Contributions are welcome! Feel free to open issues and submit pull requests.
