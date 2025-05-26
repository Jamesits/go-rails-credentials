package credentials

import (
	"bytes"
	"fmt"
)

// Implements Ruby object marshal protocol.
// https://docs.ruby-lang.org/en/2.1.0/marshal_rdoc.html
// https://github.com/hyrious/marshal
// https://unfit-for.work/posts/2023/rails-go-shared-credentials/

// UnmarshalSingleString extracts a single string from a Ruby marshalled object.
// The string must be the first item. Everything else is discarded.
func UnmarshalSingleString(marshalledObject []byte) (string, error) {
	// version
	if (marshalledObject[0] != 0x04) || (marshalledObject[1] != 0x08) {
		return "", fmt.Errorf("unknown marshal format %02x%02x", marshalledObject[0], marshalledObject[1])
	}

	// type
	if marshalledObject[2] != 0x22 {
		return "", fmt.Errorf("unknown object type: %02x", marshalledObject[2])
	}

	// length
	var length int
	var start int
	switch marshalledObject[3] {
	case 0x00, 0xfc, 0xfd, 0xfe, 0xff:
		return "", fmt.Errorf("unsupported object length: %02x", marshalledObject[3])

	case 0x01:
		length = int(marshalledObject[4])
		start = 5

	case 0x02:
		length = int(marshalledObject[4]) + int(marshalledObject[5])*256
		start = 6

	case 0x03:
		length = int(marshalledObject[4]) + int(marshalledObject[5])*256 + int(marshalledObject[6])*65536
		start = 7

	case 0x04:
		length = int(marshalledObject[4]) + int(marshalledObject[5])*256 + int(marshalledObject[6])*65536 + int(marshalledObject[7])*16777216
		start = 8

	default:
		length = int(marshalledObject[3]) - 5
		start = 4
	}
	if len(marshalledObject) < length+start {
		return "", fmt.Errorf("length validation failed, requires %d, has %d", length+start, len(marshalledObject))
	}

	return string(marshalledObject[start : start+length]), nil
}

// MarshalSingleString converts a string into Ruby marshal format.
func MarshalSingleString(source string) ([]byte, error) {
	b := bytes.Buffer{}
	b.Write([]byte{
		0x04, 0x08, // version
		0x22, // type: string
	})

	length := len(source)
	if length < 123 {
		b.WriteByte(byte(length + 5))
	} else if length < 256 {
		b.WriteByte(0x01)
		b.WriteByte(byte(length))
	} else if length < 65536 {
		b.WriteByte(0x02)
		b.WriteByte(byte(length))
		b.WriteByte(byte(length / 256))
	} else if length < 16777216 {
		b.WriteByte(0x03)
		b.WriteByte(byte(length))
		b.WriteByte(byte(length / 256))
		b.WriteByte(byte(length / 65536))
	} else {
		b.WriteByte(0x04)
		b.WriteByte(byte(length))
		b.WriteByte(byte(length / 256))
		b.WriteByte(byte(length / 65536))
		b.WriteByte(byte(length / 16777216))
	}

	b.WriteString(source)

	return b.Bytes(), nil
}
