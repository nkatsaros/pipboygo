package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"sync/atomic"
)

var pong = []byte{0, 0, 0, 0, 0}

type Command struct {
	Type int
	Args []interface{}
}

type command struct {
	Type int           `json:"type"`
	Args []interface{} `json:"args"`
	ID   int           `json:"id"`
}

type response struct {
	ID      int    `json:"id"`
	Allowed bool   `json:"allowed"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (c command) MarshalBinary() (data []byte, err error) {
	data, err = json.Marshal(c)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err = binary.Write(buf, binary.LittleEndian, uint32(len(data))); err != nil {
		return nil, err
	}
	if err = buf.WriteByte(byte(ChannelRequest)); err != nil {
		return nil, err
	}
	if _, err = buf.Write(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Encoder struct {
	w  io.Writer
	id int64
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Encode(cmd Command) (id int, err error) {
	nextID := atomic.AddInt64(&enc.id, 1)
	id = int(nextID)

	data, err := command{
		Type: cmd.Type,
		Args: cmd.Args,
		ID:   id,
	}.MarshalBinary()
	if err != nil {
		return id, err
	}

	_, err = enc.w.Write(data)
	return id, err
}

func (enc *Encoder) EncodePong() error {
	_, err := enc.w.Write(pong)
	return err
}
