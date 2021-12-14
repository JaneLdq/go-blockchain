package p2p

const CTLMSG_LEN = 12

func commandToBytes(command string) []byte {
	var bytes [CTLMSG_LEN]byte

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
