package resp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	ls := strings.Repeat("abc", 256)

	tests := []struct {
		data   []byte
		tokens []string
		err    bool
	}{
		{[]byte{}, nil, true},
		{
			data:   []byte("*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\ntest1\r\n"),
			tokens: []string{"*3", "$3", "SET", "$4", "name", "$5", "test1"},
			err:    false,
		},
		{
			data:   []byte("*3\r\n$3\r\nSET\r\n$8\r\njsonData\r\n$12\r\n{\"test\":123}\r\n"),
			tokens: []string{"*3", "$3", "SET", "$8", "jsonData", "$12", "{\"test\":123}"},
			err:    false,
		},
		{
			data:   []byte(fmt.Sprintf("*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$%d\r\n%s\r\n", len(ls), ls)),
			tokens: []string{"*3", "$3", "SET", "$4", "name", fmt.Sprintf("$%d", len(ls)), ls},
			err:    false,
		},
		{
			data:   []byte(fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(ls), ls, len(ls), ls)),
			tokens: []string{"*3", "$3", "SET", fmt.Sprintf("$%d", len(ls)), ls, fmt.Sprintf("$%d", len(ls)), ls},
			err:    false,
		},
	}

	for _, test := range tests {
		tokens, err := Decode(test.data)
		if test.err {
			assert.NotNil(t, err)
			continue
		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, test.tokens, tokens)
	}
}

func TestRESPParser_Parse(t *testing.T) {
	tests := []struct {
		tokens []string
		err    bool
		name   string
		params []string
	}{
		{tokens: []string{"*3", "$3", "set", "$4", "name", "$5", "test1"}, err: false, name: "SET", params: []string{"name", "test1"}},
		{tokens: []string{"*3", "$3", "set", "$4", "name", "$6", "test1"}, err: true, name: "SET", params: []string{"name"}},
		{tokens: []string{"*2", "$3", "get", "$4", "name"}, err: false, name: "GET", params: []string{"name"}},
		{tokens: []string{"*2", "$3", "get", "$10", "name"}, err: true, name: "GET", params: []string{}},
		{tokens: []string{"*2", "$3", "del", "$4", "name"}, err: false, name: "DEL", params: []string{"name"}},
		{tokens: []string{"*3", "$3", "del", "$4", "name", "$5", "name2"}, err: false, name: "DEL", params: []string{"name", "name2"}},
	}

	for _, test := range tests {
		rp := RESPParser{tokens: test.tokens}
		cmd, err := rp.Parse()

		if test.err {
			assert.NotNil(t, err)
			continue
		} else {
			assert.Nil(t, err)
		}

		assert.Equal(t, test.name, cmd.Name)
		assert.Equal(t, test.params, cmd.Params)
	}
}

func TestDecodeAndParse(t *testing.T) {
	ls := strings.Repeat("abc", 256)

	tests := []struct {
		data   []byte
		name   string
		params []string
		err    bool
	}{
		{
			data:   []byte("*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\ntest1\r\n"),
			params: []string{"name", "test1"},
			name:   "SET",
			err:    false,
		},
		{
			data:   []byte("*3\r\n$3\r\nSET\r\n$8\r\njsonData\r\n$12\r\n{\"test\":123}\r\n"),
			params: []string{"jsonData", "{\"test\":123}"},
			name:   "SET",
			err:    false,
		},
		{
			data:   []byte(fmt.Sprintf("*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$%d\r\n%s\r\n", len(ls), ls)),
			params: []string{"name", ls},
			name:   "SET",
			err:    false,
		},
		{
			data:   []byte(fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(ls), ls, len(ls), ls)),
			params: []string{ls, ls},
			name:   "SET",
			err:    false,
		},
	}

	for _, test := range tests {
		tokens, err := Decode(test.data)
		assert.Nil(t, err)

		rp := New(tokens)
		cmd, err := rp.Parse()

		assert.Nil(t, err)
		assert.Equal(t, test.name, cmd.Name)
		assert.Equal(t, test.params, cmd.Params)
	}
}

func BenchmarkParse(b *testing.B) {
	tokens := []string{"*3", "$3", "SET", "$8", "jsonData", "$12", "{\"test\":123}"}
	rp := New(tokens)

	for n := 0; n < b.N; n++ {
		rp.Parse()
		rp.current = 0
	}
}
