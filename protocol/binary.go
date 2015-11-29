package protocol

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type PipboyDictionary map[uint32]interface{}

var ErrInvalidBinaryProtocol = errors.New("invalid binary protocol")

func (d PipboyDictionary) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}

	for k, v := range d {
		m[fmt.Sprintf("%d", k)] = v
	}

	return json.Marshal(m)
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

// func resolve(memory map[uint32]interface{}, value uint32) (r interface{}, err error) {
//   v, ok := memory[value]
//   if !ok {
//     return nil, fmt.Errorf("invalid value")
//     // return "", nil
//   }

//   switch t := v.(type) {
//   case bool, uint8, int8, uint32, int32, float32, string:
//     return t, nil
//   case []uint32:
//     res := []interface{}{}
//     for _, thing := range t {
//       ires, err := resolve(memory, thing)
//       if err != nil {
//         return nil, err
//       }
//       res = append(res, ires)
//     }
//     return res, nil
//   case PipboyDictionary:
//     res := map[string]interface{}{}
//     for location, name := range t {
//       ires, err := resolve(memory, location)
//       if err != nil {
//         return nil, err
//       }
//       nameStr, ok := name.(string)
//       if !ok {
//         return nil, fmt.Errorf("invalid value")
//       }
//       res[nameStr] = ires
//     }
//     return res, nil
//   default:
//     return nil, fmt.Errorf("invalid type")
//   }
// }

func UnmarshalBinary(r io.Reader) (memory PipboyDictionary, err error) {
	currentOffset := 0

	br := bufio.NewReader(r)

	memory = PipboyDictionary{}

loop:
	for {
		var command byte
		var addr uint32

		err := binary.Read(br, binary.LittleEndian, &command)
		switch {
		case err == io.EOF:
			break loop
		case err != nil:
			return nil, err
		}
		currentOffset += 1

		if err = binary.Read(br, binary.LittleEndian, &addr); err != nil {
			return nil, err
		}
		currentOffset += 4

		switch command {
		case 0:
			// flag
			var flagVal byte
			var flag bool
			if err := binary.Read(br, binary.LittleEndian, &flagVal); err != nil {
				return nil, err
			}
			currentOffset += 1
			if flagVal == 0 {
				flag = false
			} else {
				flag = true
			}

			memory[addr] = flag
		case 1:
			// value
			var value int8
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			currentOffset += 1

			memory[addr] = value
		case 2:
			// value
			var value uint8
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			currentOffset += 1

			memory[addr] = value
		case 3:
			// value
			var value int32
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			currentOffset += 4

			memory[addr] = value
		case 4:
			// value
			var value uint32
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			currentOffset += 4

			memory[addr] = value
		case 5:
			// value
			var value float32
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			currentOffset += 4

			memory[addr] = value
		case 6:
			// string
			str, err := readString(br)
			if err != nil {
				return nil, err
			}
			currentOffset += len(str) + 1

			memory[addr] = str
		case 7:
			// array
			var elements uint16

			if err := binary.Read(br, binary.LittleEndian, &elements); err != nil {
				return nil, err
			}
			currentOffset += 2

			subaddrs := []uint32{}
			for i := 0; i < int(elements); i++ {
				var subaddr uint32
				if err := binary.Read(br, binary.LittleEndian, &subaddr); err != nil {
					return nil, err
				}
				currentOffset += 4

				subaddrs = append(subaddrs, subaddr)
			}
			memory[addr] = subaddrs
		case 8:
			// dictionary
			var elements uint16
			var nuls uint16

			if err := binary.Read(br, binary.LittleEndian, &elements); err != nil {
				return nil, err
			}
			currentOffset += 2

			dictionary := map[string]uint32{}
			for i := 0; i < int(elements); i++ {
				var subaddr uint32
				if err := binary.Read(br, binary.LittleEndian, &subaddr); err != nil {
					return nil, err
				}
				currentOffset += 4

				str, err := readString(br)
				if err != nil {
					return nil, err
				}
				currentOffset += len(str) + 1

				dictionary[str] = subaddr
			}

			// dictionaries have a two NUL bytes at the end
			if err := binary.Read(br, binary.LittleEndian, &nuls); err != nil {
				return nil, err
			}
			if nuls != 0x0000 {
				return nil, ErrInvalidBinaryProtocol
			}
			currentOffset += 2
			memory[addr] = dictionary
		default:
			return nil, ErrInvalidBinaryProtocol
		}
	}

	return memory, nil
}
