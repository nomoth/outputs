package ipc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
)

type msgType uint32

const (
	RUN_COMMAND msgType = iota
	GET_WORKSPACES
	SUBSCRIBE
	GET_OUTPUTS
	GET_TREE
	GET_MARKS
	GET_BAR_CONFIG
	GET_VERSION
	GET_BINDING_MODES
	GET_CONFIG
	SEND_TICK
	SYNC
)

type Connection struct {
	net.Conn
}

type Mode struct {
	Width   int
	Height  int
	Refresh int
}

type Output struct {
	Name        string
	Make        string
	Model       string
	Serial      string
	Active      bool
	Scale       float32
	Modes       []*Mode
	CurrentMode *Mode `json:"current_mode"`
}

func (c *Connection) Run(cmd string) error {
	j, err := c.send(RUN_COMMAND, cmd)
	if err != nil {
		return err
	}
	var status []struct {
		Success bool
		Error   string
	}
	err = json.Unmarshal(j, &status)
	if err != nil {
		return err
	}
	for _, s := range status {
		if !s.Success {
			return errors.New(s.Error)
		}
	}
	return nil
}

func (c *Connection) GetOutputs() ([]*Output, error) {
	j, err := c.send(GET_OUTPUTS, "")
	if err != nil {
		return nil, err
	}
	var outputs []*Output
	err = json.Unmarshal(j, &outputs)
	if err != nil {
		return nil, err
	}
	return outputs, nil
}

func NewConnection() *Connection {
	return &Connection{
		Conn: getConnection(getSocketPath()),
	}
}

var magic = [6]byte{'i', '3', '-', 'i', 'p', 'c'}

func (c *Connection) send(t msgType, cmd string) ([]byte, error) {
	h := &struct {
		Magic  [6]byte
		Length uint32
		Type   msgType
	}{
		Magic:  magic,
		Length: uint32(len(cmd)),
		Type:   t,
	}
	err := binary.Write(c, NativeByteOrder, h)
	if err != nil {
		return nil, err
	}
	_, err = c.Write([]byte(cmd))
	if err != nil {
		return nil, err
	}
	err = binary.Read(c, NativeByteOrder, h)
	if err != nil {
		return nil, err
	}
	b := bytes.Buffer{}
	_, err = io.CopyN(&b, c, int64(h.Length))
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func getConnection(address string) net.Conn {
	c, err := net.Dial("unix", address)
	if err != nil {
		log.Fatalf("impossible to connect: %s", err)
	}
	return c
}

func getSocketPath() string {
	c := exec.Command("sway", "--get-socketpath")
	var out bytes.Buffer
	c.Stdout = &out
	err := c.Run()
	if err != nil {
		log.Fatalf("impossible to get the socket path: %s", err)
	}
	return strings.TrimSpace(out.String())
}
