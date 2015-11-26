package protocol

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
)

const (
	ChannelHeartbeat = 0
	ChannelGameInfo  = 1
	ChannelBusy      = 2
	ChannelData      = 3
	ChannelLocalMap  = 4
	ChannelRequest   = 5
	ChannelResponse  = 6
)

// {"lang": "en", "version": "1.1.30.0"}
type GameInfo struct {
	Language string `json:"lang"`
	Version  string `json:"version"`
}

type CommandResponse struct {
	ID      int    `json:"id"`
	Allowed bool   `json:"allowed"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type LocalMap struct {
	Image   image.Image
	Extents Extents
}

type Decoder struct {
	r   io.Reader
	br  *bufio.Reader
	err error
}

type UnknownChannelError struct {
	Channel int
	Data    []byte
}

func (e *UnknownChannelError) Error() string {
	return fmt.Sprintf("channel %d not implemented", e.Channel)
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:  r,
		br: bufio.NewReader(r),
	}
}

func (dec *Decoder) Decode() (int, interface{}, error) {
	if dec.err != nil {
		return 0, nil, dec.err
	}

	var size uint32
	var chanbyte byte

	dec.read(&size)
	dec.read(&chanbyte)

	if dec.err != nil {
		return 0, nil, dec.err
	}

	channel := int(chanbyte)

	lr := io.LimitReader(dec.br, int64(size))

	// drain the rest of the channel data in case our handlers suck
	defer io.Copy(ioutil.Discard, lr)

	switch channel {
	case ChannelHeartbeat:
		return channel, nil, nil
	case ChannelGameInfo:
		var info GameInfo
		if err := json.NewDecoder(lr).Decode(&info); err != nil {
			return 0, nil, err
		}
		return channel, info, nil
	case ChannelBusy:
		return channel, nil, nil
	case ChannelData:
		memory, err := UnmarshalBinary(lr)
		if err != nil {
			return 0, nil, err
		}
		return channel, memory, nil
	case ChannelLocalMap:
		image, extents, err := UnmarshalMap(lr, int(size))
		if err != nil {
			return 0, nil, err
		}
		return channel, LocalMap{image, extents}, nil
	case ChannelResponse:
		var response CommandResponse
		if err := json.NewDecoder(lr).Decode(&response); err != nil {
			return 0, nil, err
		}
		return channel, response, nil
	}

	// read all data from the unknown channel
	data, err := ioutil.ReadAll(lr)
	if err != nil {
		return 0, nil, err
	}

	return 0, nil, &UnknownChannelError{channel, data}
}

func (dec *Decoder) read(data interface{}) {
	if dec.err != nil {
		return
	}
	dec.err = binary.Read(dec.br, binary.LittleEndian, data)
}

func (dec *Decoder) Buffered() io.Reader {
	return io.LimitReader(dec.br, int64(dec.br.Buffered()))
}
