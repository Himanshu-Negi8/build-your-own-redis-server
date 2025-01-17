package parser

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Himanshu-Negi8/build-your-own-redis-server/types"
	"io"
	"strconv"
	"strings"
)

/*
	Parse, parse command input.
	This is built on top of the RESP protocol, which is a serialization protocol used by Redis.
	Examples of RESP protocol handling can be found in the Redis documentation: https://redis.io/topics/protocol
	The function reads the first byte from the reader and based on the prefix, it reads the next bytes accordingly.
	The function returns the parsed data, the type of the data, and an error if any.
*/

func Parse(r io.Reader) (interface{}, types.RESPType, error) {
	/*
		The `bufio.Reader` is used to read data from a buffered stream, which is essentially a stream of bytes.
		Once the bytes are processed, they can be converted to the desired type. As you read from a `bufio.Reader`,
		the bytes are consumed from the underlying stream, meaning they are no longer available for subsequent reads.
	*/
	reader := bufio.NewReader(r)

	// Read the first byte
	prefix, err := reader.ReadByte()
	if err != nil {
		return nil, "", err
	}

	switch prefix {
	// If the prefix is '+', it is a Simple String.
	case '+': // Simple String
		line, err := readLine(reader)
		if err != nil {
			return nil, "", err
		}
		return line, types.RESPTypeSimpleString, nil
	// If the prefix is '-', it is an Error.
	case '-': // Error
		line, err := readLine(reader)
		if err != nil {
			return nil, "", err
		}
		return fmt.Errorf("redis error: %s", line), types.RESPTypeError, nil
	case ':': // Integer
		line, err := readLine(reader)
		if err != nil {
			return nil, "", err
		}
		intValue, err := strconv.Atoi(line)
		if err != nil {
			return nil, "", err
		}
		return intValue, types.RESPTypeInteger, nil
	// If the prefix is '$', it is a Bulk String.
	// Example of a Bulk String: $6\r\nfoobar\r\n
	case '$': // Bulk String
		length, err := readLength(reader)
		if err != nil {
			return nil, "", err
		}
		if length == -1 {
			return nil, types.RESPTypeBulkString, nil // Null bulk string
		}
		data := make([]byte, length)
		if _, err := io.ReadFull(reader, data); err != nil {
			return nil, "", err
		}
		// Consume the trailing \r\n
		if _, err := readLine(reader); err != nil {
			return nil, "", err
		}
		return string(data), types.RESPTypeBulkString, nil
	// If the prefix is '*', it is an Array.
	// Example of an Array: *2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n
	// This is an array with two elements: "foo" and "bar".
	case '*': // Array
		length, err := readLength(reader)
		if err != nil {
			return nil, "", err
		}
		if length == -1 {
			return nil, types.RESPTypeArray, nil // Null array
		}
		elements := make([]interface{}, length)
		for i := 0; i < length; i++ {
			// Since it's array we need to recursively call ParseRESP to parse each element.
			elem, _, err := Parse(reader)
			if err != nil {
				return nil, "", err
			}
			elements[i] = elem
		}
		return elements, types.RESPTypeArray, nil
	default:
		return nil, "", errors.New("unknown prefix: " + string(prefix))
	}
}

// Helper to read a line ending with \r\n
// readLine reads bytes from the reader until it encounters a '\n' character.
func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(line, "\r\n"), nil
}

// Helper to parse length or integer
// readLength calls readLine and changes the first byte into an integer.
func readLength(reader *bufio.Reader) (int, error) {
	line, err := readLine(reader)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(line)
}
