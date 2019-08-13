package erds

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/tidwall/redcon"
)

type aofReader struct {
	input *bufio.Reader
}

func newAofReader(reader io.Reader) aofReader {
	input := bufio.NewReader(reader)
	return aofReader{input: input}
}

func (reader aofReader) readlineAndCount() (line []byte, count int, err error) {
	line, err = reader.input.ReadBytes('\n')
	if err != nil {
		return
	}
	str := string(line)
	if string(str[0]) == "*" {
		count, _ = strconv.Atoi(str[1:2])
	} else {
		count = 0
	}
	return
}

func (reader aofReader) readParameter() (para []byte) {
	// read parameter length
	para, _ = reader.input.ReadBytes('\n')
	str := string(para)
	if string(str[0]) != "$" {
		panic("Corrupt File: Element is not parameter length")
	}
	// read parameter
	next, _ := reader.input.ReadBytes('\n')
	para = append(para, next...)
	return
}

func (reader aofReader) readOneRESP() (resp []byte, err error) {
	line, count, err := reader.readlineAndCount()
	if err != nil {
		return
	}
	resp = append(resp, line...)
	for ; count > 0; count-- {
		resp = append(resp, reader.readParameter()...)
	}
	return
}

const aofStatusOn = 0
const aofStatusOff = 1

type aof struct {
	status int
	buffer *bytes.Buffer
}

func newAof() *aof {
	return &aof{
		status: aofStatusOn,
		buffer: bytes.NewBuffer([]byte{}),
	}
}

func (aof *aof) isTurnOn() bool {
	return aof.status == aofStatusOn
}

func (aof *aof) propagateAof(key string, cmd redcon.Command) {
	switch key {
	case "set":
		if len(cmd.Args) != 3 {
			return
		}
	case "del":
		if len(cmd.Args) != 2 {
			return
		}
	default:
		return
	}
	aof.buffer.Write(cmd.Raw)
}

func (aof *aof) saveAofFile() {
	if !aof.isTurnOn() {
		return
	}

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("erds.aof", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(aof.buffer.Bytes()); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	// reset buffer
	aof.buffer.Reset()
}
