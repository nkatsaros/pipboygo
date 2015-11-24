package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func readString(r *bufio.Reader) (str string, err error) {
	line, err := r.ReadBytes(0x00)
	if err != nil {
		return "", err
	}
	// remove the NUL byte
	line = line[:len(line)-1]
	return string(line), nil
}

func resolve(memory map[uint32]interface{}, value uint32) (r interface{}, err error) {
	v, ok := memory[value]
	if !ok {
		return nil, fmt.Errorf("invalid value")
	}

	switch t := v.(type) {
	case bool, uint8, int8, uint32, int32, float32, string:
		return t, nil
	case []uint32:
		res := []interface{}{}
		for _, thing := range t {
			ires, err := resolve(memory, thing)
			if err != nil {
				return nil, err
			}
			res = append(res, ires)
		}
		return res, nil
	case map[uint32]string:
		res := map[string]interface{}{}
		for location, name := range t {
			ires, err := resolve(memory, location)
			if err != nil {
				return nil, err
			}
			res[name] = ires
		}
		return res, nil
	default:
		return nil, fmt.Errorf("invalid type")
	}
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	currentOffset := 0

	if len(os.Args) >= 3 {
		offset, err := strconv.ParseInt(os.Args[2], 16, 32)
		if err != nil {
			log.Fatalln(err)
		}
		currentOffset = int(offset)
		log.Println("seeking to", offset)
		_, err = f.Seek(offset, os.SEEK_SET)
		if err != nil {
			log.Fatalln(err)
		}
	}

	br := bufio.NewReader(f)

	memory := map[uint32]interface{}{}

loop:
	for {
		var command byte
		var addr uint32

		err := binary.Read(br, binary.LittleEndian, &command)
		switch {
		case err == io.EOF:
			break loop
		case err != nil:
			log.Fatalln(err)
		}
		log.Printf("0x%X", currentOffset)
		currentOffset += 1

		if err = binary.Read(br, binary.LittleEndian, &addr); err != nil {
			fmt.Println(addr, err)
			log.Println(command, "here2!")
			log.Fatalln(err)
		}
		currentOffset += 4

		switch command {
		case 0:
			// flag
			var flagVal byte
			var flag bool
			if err := binary.Read(br, binary.LittleEndian, &flagVal); err != nil {
				log.Fatalln(err)
			}
			currentOffset += 1
			if flagVal == 0 {
				flag = false
			} else {
				flag = true
			}
			log.Println("Flag", addr, flag)
			memory[addr] = flag
		case 1:
			// value
			var value int8
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				log.Fatalln(err)
			}
			currentOffset += 1
			log.Println(command, addr, value)
			memory[addr] = value
		case 2:
			// value
			var value uint8
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				log.Fatalln(err)
			}
			currentOffset += 1
			log.Println(command, addr, value)
			memory[addr] = value
		case 3:
			// value
			var value int32
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				log.Fatalln(err)
			}
			currentOffset += 4
			log.Println(command, addr, value)
			memory[addr] = value
		case 4:
			// value
			var value uint32
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				log.Fatalln(err)
			}
			currentOffset += 4
			log.Println(command, addr, value)
			memory[addr] = value
		case 5:
			// value
			var value float32
			if err := binary.Read(br, binary.LittleEndian, &value); err != nil {
				log.Fatalln(err)
			}
			currentOffset += 4
			log.Println(command, addr, value)
			memory[addr] = value
		case 6:
			// string
			str, err := readString(br)
			if err != nil {
				log.Fatalln(err)
			}
			currentOffset += len(str) + 1
			log.Println(command, addr, str)
			memory[addr] = str
		case 7:
			// array
			var elements uint16

			if err := binary.Read(br, binary.LittleEndian, &elements); err != nil {
				log.Fatalln(err)
			}
			currentOffset += 2

			log.Println(command, addr, elements, "elements")

			subaddrs := []uint32{}
			for i := 0; i < int(elements); i++ {
				var subaddr uint32
				if err := binary.Read(br, binary.LittleEndian, &subaddr); err != nil {
					log.Fatalln(err)
				}
				currentOffset += 4
				log.Println("subaddr", subaddr)
				subaddrs = append(subaddrs, subaddr)
			}
			memory[addr] = subaddrs
		case 8:
			// dictionary
			var elements uint16
			var nuls uint16

			if err := binary.Read(br, binary.LittleEndian, &elements); err != nil {
				log.Fatalln(err)
			}
			currentOffset += 2

			log.Println(command, addr, elements, "elements")

			dictionary := map[uint32]string{}
			for i := 0; i < int(elements); i++ {
				var subaddr uint32
				if err := binary.Read(br, binary.LittleEndian, &subaddr); err != nil {
					log.Fatalln(err)
				}
				currentOffset += 4

				str, err := readString(br)
				if err != nil {
					log.Fatalln(err)
				}
				currentOffset += len(str) + 1
				log.Println("subaddr", subaddr, str)
				dictionary[subaddr] = str
			}

			// dictionaries have a two NUL bytes at the end
			if err := binary.Read(br, binary.LittleEndian, &nuls); err != nil {
				log.Fatalln(err)
			}
			if nuls != 0x0000 {
				log.Fatalln("should be NUL")
			}
			currentOffset += 2
			memory[addr] = dictionary
		default:
			log.Fatalln("unsupported command", command)
		}
	}

	blah, err := resolve(memory, 0)
	if err != nil {
		log.Fatalln(err)
	}
	data, err := json.Marshal(blah)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(data))
}
