/*
* jprot: jerry protocol
*
* NUMBER		{+-}?{0..9}+		Max. length   9 digits + 1 for sign
* TEXTSTRING	{Unicode char}*		Max. length 256 bytes
* MESSAGE		{<NUMBER> | <TEXTSTRING>}+
*
* Encoding uses LENGTH-CONTENT encoding:
* LENGTH	1 byte, value is the number of content bytes
* CONTENT	0 or more content bytes
*
* Example:
* "hello"	>	[]byte{ 5 , h , e , l , l , o}
* 42		>	[]byte{ 2, '4', '2'}
*
 */
package jprot

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

const (
	TagNumber byte = 0x01
	TagString byte = 0x02
)

type myreader interface {
	io.Reader
	io.ByteReader
}

func Marshal(msg []interface{}) ([]byte, error) {
	var result bytes.Buffer
	result.Grow(1024)

	for _, m := range msg {
		switch m.(type) {
		case string:
			s, err := marshalString(m.(string))
			if err != nil {
				return nil, err
			}
			result.WriteByte(TagString)
			result.Write(s)
		case int:
			i, err := marshalInt(m.(int))
			if err != nil {
				return nil, err
			}
			result.WriteByte(TagNumber)
			result.Write(i)
		default:
			return nil, fmt.Errorf("Marshal: Unknown type in message: %v", m)
		}
	}

	b := result.Bytes()
	return b[:len(b)], nil
}

// marshalString: Marshals a string into jprot byte sequence
// Max String Length: 255
func marshalString(v string) ([]byte, error) {
	if len(v) > 255 {
		return nil, fmt.Errorf("Marshal: string length too high! %d", len(v))
	}

	result := make([]byte, len(v)+1)
	result[0] = byte(len(v))

	copy(result[1:], v)

	return result, nil
}

func marshalInt(v int) ([]byte, error) {
	str := strconv.Itoa(v)
	if len(str) > 10 {
		return nil, fmt.Errorf("Marshal: Number too large! %d digits", len(str))
	}

	result := make([]byte, len(str)+1)
	result[0] = byte(len(str))
	copy(result[1:], []byte(str))

	return result, nil
}

func Unmarshal(r myreader) ([]interface{}, error) {
	result := []interface{}{}

	for {
		tag, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch tag {
		case TagNumber:
			num, err := unmarshalInt(r)
			if err != nil {
				return nil, err
			}
			result = append(result, num)
		case TagString:
			str, err := unmarshalString(r)
			if err != nil {
				return nil, err
			}
			result = append(result, str)
		default:
			return nil, fmt.Errorf("Unmarshal: Unknown Tag %x", tag)
		}
	}

	return result, nil
}

func unmarshalString(r myreader) (string, error) {
	l, err := r.ReadByte()
	if err != nil {
		return "", err
	}

	result := make([]byte, l)

	if _, err := r.Read(result); err != nil {
		return "", err
	}

	return string(result), nil
}

// unmarshalInt: Unmarshals a single value
// Max Length: 10
func unmarshalInt(r myreader) (int, error) {
	l, err := r.ReadByte()
	if err != nil {
		return 0, err
	}

	if l > 10 {
		return 0, fmt.Errorf("Unmarshal: Number content length too high! %d", l)
	}

	p := make([]byte, l)

	_, err = r.Read(p)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(p))
}
