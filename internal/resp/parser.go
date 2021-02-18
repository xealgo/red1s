package resp

import (
	"fmt"
	"strconv"
	"strings"
)

// Command basic structure to hold parsed RESP commands.
type Command struct {
	Name string

	// This could be improved later by creating a Param struct
	// which could hold the data types for each parameter.
	Params []string
}

// RESPParser implements methods for parsing the RESP data.
type RESPParser struct {
	current int
	tokens  []string
}

// New returns a new RESPParser
func New(tokens []string) *RESPParser {
	return &RESPParser{tokens: tokens}
}

// Parse parses the RESP command.
func (r *RESPParser) Parse() (*Command, error) {
	token := r.next()
	if token == nil {
		return nil, fmt.Errorf("Error decoding: %w", ErrEndOfInput)
	}

	partsCount, err := r.parseLength(*token, '*')
	if err != nil {
		return nil, fmt.Errorf("Error decoding: %w", err)
	}

	size := len(r.tokens)

	// The number of tokens should match the partsCount * 2 since there is an preceding
	// size token for each command part: $2 set %3 num %1 4 (6 tokens, 3 parts).
	if partsCount == 0 || size-1 != partsCount*2 {
		return nil, fmt.Errorf("Error decoding: expected %d command parts, got %d", partsCount*2, size)
	}

	cmd, err := r.parseCommand()
	if err != nil {
		return nil, fmt.Errorf("Error decoding: %w", err)
	}
	return cmd, nil
}

func (r *RESPParser) next() *string {
	if len(r.tokens) == r.current {
		return nil
	}

	token := r.tokens[r.current]
	r.current++

	return &token
}

func (r *RESPParser) parseLength(token string, delimeter byte) (int, error) {
	if token[0] != delimeter {
		return 0, fmt.Errorf("Invalid syntax: Expected '%s', got '%s'", string(delimeter), string(token[0]))
	}

	value := token[1:]

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("Invalid syntax: Expected integer, got: '%s'", value)
	}
	return n, nil
}

func (r *RESPParser) parseCommand() (*Command, error) {
	token := r.next()
	if token == nil {
		return nil, fmt.Errorf("Error parsing command: %w", ErrEndOfInput)
	}

	n, err := r.parseLength(*token, '$')
	if err != nil {
		return nil, fmt.Errorf("Error parsing command: %w", err)
	}

	token = r.next()
	if token == nil {
		return nil, fmt.Errorf("Error parsing command: %w", ErrEndOfInput)
	}

	name := *token

	if len(name) != n {
		return nil, fmt.Errorf("Error parsing command. Expected length %d, got %d", n, len(name))
	}

	cmd := Command{Name: strings.ToUpper(name)}

	cmd.Params, err = r.parseParams()
	if err != nil {
		return nil, fmt.Errorf("Error parsing command: %w", err)
	}

	return &cmd, nil
}

func (r *RESPParser) parseParams() ([]string, error) {
	var params []string

	token := r.next()
	for token != nil {
		n, err := r.parseLength(*token, '$')
		if err != nil {
			return nil, fmt.Errorf("Error parsing params: %w", err)
		}

		token = r.next()

		if len(*token) != n {
			return nil, fmt.Errorf("Error parsing params. Expected length %d, got %d, token: %s", n, len(*token), *token)
		}

		params = append(params, *token)
		token = r.next()
	}
	return params, nil
}
