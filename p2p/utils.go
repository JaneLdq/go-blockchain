package p2p

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func commandToBytes(command string) []byte {
	var bytes [P2P_CMD_LEN]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return string(command)
}

func sendData(addr string, data []byte) error {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		return err
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
	return nil
}

func getPayload(request []byte) bytes.Buffer {
	var buff bytes.Buffer
	buff.Write(request[P2P_CMD_LEN:])
	return buff
}