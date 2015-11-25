package protocol

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"
)

const mapHeaderSize = 32

type Extent struct {
	X, Y float32
}

type Extents struct {
	NW, NE, SW Extent
}

type mapReader struct {
	r   io.Reader
	err error
}

func (r *mapReader) read(data interface{}) {
	if r.err != nil {
		return
	}
	r.err = binary.Read(r.r, binary.LittleEndian, data)
}

func UnmarshalMap(r io.Reader, size int) (img image.Image, extents Extents, err error) {
	var widthu32, heightu32 uint32

	mr := &mapReader{r: r}

	mr.read(&widthu32)
	mr.read(&heightu32)
	mr.read(&extents.NW.X)
	mr.read(&extents.NW.Y)
	mr.read(&extents.NE.X)
	mr.read(&extents.NE.Y)
	mr.read(&extents.SW.X)
	mr.read(&extents.SW.Y)

	if mr.err != nil {
		return nil, extents, mr.err
	}

	width, height := int(widthu32), int(heightu32)

	// This is a fix for Fallout returning the wrong resolution.
	// Found at https://github.com/CyberShadow/csfo4/blob/master/mapfix/mapfix.d
	if width*height < size-mapHeaderSize {
		width = (size - mapHeaderSize) / height
		if size != mapHeaderSize+width*height {
			return nil, extents, errors.New("invalid map stride")
		}
	}

	gray := image.NewGray(image.Rect(0, 0, width, height))

	var pixel uint8
	for i := 0; i < width*height; i++ {
		mr.read(&pixel)

		if mr.err != nil {
			return nil, extents, mr.err
		}

		gray.SetGray(i%width, i/width, color.Gray{pixel})
	}

	return gray, extents, nil
}
