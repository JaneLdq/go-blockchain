package p2p

type CMD uint8

const CMD_LENGTH = 1

const (
	CONNECT CMD = iota
	HELLO
	MINE
	NEW_BLOCK
	REQ_CHAIN
	UPDATE_CHAIN
)

func (cmd CMD) String() string {
	return [...]string{"connect", "hello", "mine", "newblock", "reqchain", "updatechain"}[cmd]
}

func (cmd CMD) ToByteArray() []byte {
	bs := make([]byte, CMD_LENGTH)
	bs[0] = byte(cmd)
	return bs
}