package p2p

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

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

func getPayload(request []byte) []byte {
	var buff bytes.Buffer
	buff.Write(request[CMD_LENGTH:])
	return buff.Bytes()
}
