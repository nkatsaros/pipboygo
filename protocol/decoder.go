package protocol

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"io"
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

	switch channel {
	case 0:
		return channel, nil, nil
	case 1:
		var info GameInfo
		if err := json.NewDecoder(io.LimitReader(dec.br, int64(size))).Decode(&info); err != nil {
			return 0, nil, err
		}
		return channel, info, nil
	case 3:
		memory, err := UnmarshalBinary(io.LimitReader(dec.br, int64(size)))
		if err != nil {
			return 0, nil, err
		}
		return channel, memory, nil
	case 4:
		image, extents, err := UnmarshalMap(io.LimitReader(dec.br, int64(size)), int(size))
		if err != nil {
			return 0, nil, err
		}
		return channel, LocalMap{image, extents}, nil
	case 6:
		var response CommandResponse
		if err := json.NewDecoder(io.LimitReader(dec.br, int64(size))).Decode(&response); err != nil {
			return 0, nil, err
		}
		return channel, response, nil
	default:
		return 0, nil, fmt.Errorf("channel %d not implemented", channel)
	}

	return channel, nil, nil
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
