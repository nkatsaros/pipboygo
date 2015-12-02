package protocol

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"strconv"
)

type PipboyMemory struct {
	Added   map[uint32]interface{}
	Removed []uint32
}

var ErrInvalidBinaryProtocol = errors.New("invalid binary protocol")

func (d PipboyMemory) MarshalJSON() ([]byte, error) {
	r := struct {
		Added   map[string]interface{} `json:"added"`
		Removed []uint32               `json:"removed"`
	}{}
	r.Added = map[string]interface{}{}

	for k, v := range d.Added {
		r.Added[strconv.Itoa(int(k))] = v
	}

	r.Removed = d.Removed

	return json.Marshal(r)
}

func readString(r *bufio.Reader) (str string, err error) {
	line, err := r.ReadBytes(0x00)
	if err != nil {
		return "", err
	}
	// remove the NUL byte
	line = line[:len(line)-1]
	return string(line), nil
}

func UnmarshalBinary(r io.Reader) (memory PipboyMemory, err error) {
	currentOffset := 0

	br := bufio.NewReader(r)

	memory = PipboyMemory{}
	memory.Added = map[uint32]interface{}{}
	memory.Removed = []uint32{}

loop:
	for {
		var command byte
		var addr uint32

		err := binary.Read(br, binary.LittleEndian, &command)
		switch {
		case err == io.EOF:
			break loop
		case err != nil:
			return memory, err
		}
		currentOffset += 1

		if err = binary.Read(br, binary.LittleEndian, &addr); err != nil {
			return memory, err
		}
		currentOffset += 4

		switch command {
		case 0:
			// flag
			var flagVal byte
			var flag bool
			if err := binary.Read(br, binary.LittleEndian, &flagVal); err != nil {
				return memory, err
			}
			currentOffset += 1
			if flagVal == 0 {
				flag = false
			} else {
				flag = true
			}

			memory.Added[addr] = flag
		case 1:
			// value
			var value int8
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				return memory, err
			}
			currentOffset += 1

			memory.Added[addr] = value
		case 2:
			// value
			var value uint8
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				return memory, err
			}
			currentOffset += 1

			memory.Added[addr] = value
		case 3:
			// value
			var value int32
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				return memory, err
			}
			currentOffset += 4

			memory.Added[addr] = value
		case 4:
			// value
			var value uint32
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				return memory, err
			}
			currentOffset += 4

			memory.Added[addr] = value
		case 5:
			// value
			var value float32
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				return memory, err
			}
			currentOffset += 4

			memory.Added[addr] = value
		case 6:
			// string
			str, err := readString(br)
			if err != nil {
				return memory, err
			}
			currentOffset += len(str) + 1

			memory.Added[addr] = str
		case 7:
			// array
			var elements uint16

			if err := binary.Read(br, binary.LittleEndian, &elements); err != nil {
				return memory, err
			}
			currentOffset += 2

			subaddrs := []uint32{}
			for i := 0; i < int(elements); i++ {
				var subaddr uint32
				if err := binary.Read(br, binary.LittleEndian, &subaddr); err != nil {
					return memory, err
				}
				currentOffset += 4

				subaddrs = append(subaddrs, subaddr)
			}
			memory.Added[addr] = subaddrs
		case 8:
			// dictionary
			var added, removed uint16

			// number of elements added
			if err := binary.Read(br, binary.LittleEndian, &added); err != nil {
				return memory, err
			}
			currentOffset += 2

			dictionary := map[string]uint32{}
			for i := 0; i < int(added); i++ {
				var subaddr uint32
				if err := binary.Read(br, binary.LittleEndian, &subaddr); err != nil {
					return memory, err
				}
				currentOffset += 4

				str, err := readString(br)
				if err != nil {
					return memory, err
				}
				currentOffset += len(str) + 1

				dictionary[str] = subaddr
			}

			memory.Added[addr] = dictionary

			// number of elements removed
			if err := binary.Read(br, binary.LittleEndian, &removed); err != nil {
				return memory, err
			}
			currentOffset += 2

			for i := 0; i < int(removed); i++ {
				var subaddr uint32
				if err := binary.Read(br, binary.LittleEndian, &subaddr); err != nil {
					return memory, err
				}
				currentOffset += 4

				memory.Removed = append(memory.Removed, subaddr)
			}

		default:
			return memory, ErrInvalidBinaryProtocol
		}
	}

	return memory, nil
}
