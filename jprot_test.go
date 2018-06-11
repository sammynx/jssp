package jprot

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	message := []interface{}{"ID", 10856, "year", 2018}
	encoded := []byte{
		TagString, 2, 'I', 'D',
		TagNumber, 5, '1', '0', '8', '5', '6',
		TagString, 4, 'y', 'e', 'a', 'r',
		TagNumber, 4, '2', '0', '1', '8',
	}

	cases := []struct {
		val  []interface{}
		want []byte
	}{
		{message, encoded},
		{[]interface{}{}, []byte{}},
		{[]interface{}{1}, []byte{TagNumber, 1, '1'}},
		{nil, []byte{}},
	}

	for _, c := range cases {
		got, _ := Marshal(c.val)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Marshal: got %v want %v", got, c.want)
		}
	}
}

func TestMarshalString(t *testing.T) {
	buf := make([]byte, 256)
	buf[0] = 255

	maxstring := string(buf[1:])

	cases := []struct {
		val  string
		want []byte
	}{
		{"hello", []byte{5, 'h', 'e', 'l', 'l', 'o'}},
		{"", []byte{0}},
		{maxstring, buf},
	}

	for _, c := range cases {
		got, _ := marshalString(c.val)
		if !bytes.Equal(got, c.want) {
			t.Errorf("marshalString got %v want %v", got, c.want)
		}
	}

	// The Bad Case
	bigstring := string(buf)
	if _, err := marshalString(bigstring); err == nil {
		t.Error("marshalString: No Error with a too large string!")
	}
}

func TestMarshalInt(t *testing.T) {
	cases := []struct {
		val  int
		want []byte
	}{
		{0, []byte{1, 48}},
		{-1, []byte{2, 45, 49}},
		{9999999999, []byte{10, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57}},
		{-999999999, []byte{10, 45, 57, 57, 57, 57, 57, 57, 57, 57, 57}},
	}

	for _, c := range cases {
		got, _ := marshalInt(c.val)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("marshalInt: got %v want %v", got, c.want)
		}
	}

	// The Bad Case
	if _, err := marshalInt(99999999999); err == nil {
		t.Error("marshalInt: No Error on number with more than max digits")
	}
}

func TestUnMarshal(t *testing.T) {
	m := []byte{TagNumber, 1, 49, TagString, 2, 'h', 'i', TagNumber, 1, 50}
	w := []interface{}{1, "hi", 2}

	cases := []struct {
		msg  []byte
		want []interface{}
	}{
		{m, w},
		{[]byte{}, []interface{}{}},
		{m[:3], []interface{}{1}},
	}

	for _, c := range cases {
		r := bytes.NewReader(c.msg)
		got, _ := Unmarshal(r)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Unmarshal: got %v want %v", got, c.want)
		}
	}
}

func TestUnmarshalIntGood(t *testing.T) {
	cases := []struct {
		val  []byte
		want int
	}{
		{[]byte{1, 48}, 0},
		{[]byte{2, 45, 49}, -1},
		{[]byte{10, 45, 57, 57, 57, 57, 57, 57, 57, 57, 57}, -999999999},
		{[]byte{10, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57}, 9999999999},
	}

	for _, c := range cases {
		r := bytes.NewReader(c.val)
		got, _ := unmarshalInt(r)
		if got != c.want {
			t.Errorf("unmarshalInt got %d want %d", got, c.want)
		}
	}
}

func TestUnmarshalIntBad(t *testing.T) {
	cases := []struct {
		val []byte
	}{
		{[]byte{}},
		{[]byte{1}},
		{[]byte{11, 45, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57}},
	}

	for _, c := range cases {
		r := bytes.NewReader(c.val)
		_, err := unmarshalInt(r)
		if err == nil {
			t.Errorf("unmarshalInt: No error on wrong input %v", c.val)
		}
	}
}

func TestUnmarshalStringGood(t *testing.T) {
	buf := make([]byte, 256)
	buf[0] = 255

	cases := []struct {
		val  []byte
		want string
	}{
		{[]byte{1, 65}, "A"},
		{[]byte{0}, ""},
		{buf, string(buf[1:])},
	}

	for _, c := range cases {
		r := bytes.NewReader(c.val)
		got, _ := unmarshalString(r)
		if got != c.want {
			t.Errorf("unmarshalString: got %s want %s", got, c.want)
		}
	}
}

func TestUnmarshalStringBad(t *testing.T) {

	cases := []struct {
		val []byte
	}{
		{[]byte{2, 65}},
		{[]byte{}},
		{nil},
	}

	for _, c := range cases {
		r := bytes.NewReader(c.val)
		_, err := unmarshalString(r)
		if err == nil {
			t.Errorf("unmarshalString: Got no error on wrong input! %v", c.val)
		}
	}
}

func TestMarshalUnmarshal(t *testing.T) {

	message := []interface{}{"ID", 10856, "year", 2018, "month", 6}

	m, err := Marshal(message)
	if err != nil {
		t.Errorf("MarshalUnmarshal: %s", err)
	}
	r := bytes.NewReader(m)

	u, err := Unmarshal(r)
	if err != nil {
		t.Errorf("MarshalUnmarshal: %s", err)
	}

	if !reflect.DeepEqual(u, message) {
		fmt.Println("MarshalUnmarshal: ", err)
		t.Errorf("Marshal->Unmarshal Error. got %v want %v", m, message)
	}
}
