package internal

import (
	"strconv"

	"github.com/xealgo/red1s/internal/resp"
	"github.com/xealgo/red1s/internal/store"
)

// Supported commands.
const (
	CmdSet = "SET"
	CmdGet = "GET"
	CmdDel = "DEL"
)

// Get grabs a value from the data store based on the provided
// command data.
func Get(ds store.DataStore, cmd *resp.Command) ([]byte, error) {
	if len(cmd.Params) == 0 {
		return []byte("-No key provided\r\n"), nil
	}

	value, err := ds.Get(cmd.Params[0])
	if err != nil {
		return []byte("$-1\r\n"), nil
	}
	return []byte("+" + *value + "\r\n"), nil
}

// Set adds/updates a value from the data store based on the provided
// command data.
func Set(ds store.DataStore, cmd *resp.Command) ([]byte, error) {
	if len(cmd.Params) < 2 {
		return []byte("-A key value pair is required\r\n"), nil
	}

	err := ds.Set(cmd.Params[0], cmd.Params[1])
	if err != nil {
		return []byte("-Key exists\r\n"), nil
	}
	return []byte("+OK\r\n"), nil
}

// Del removes entries from the data store based on the provided
// command data.
func Del(ds store.DataStore, cmd *resp.Command) ([]byte, error) {
	if len(cmd.Params) == 0 {
		return []byte("-No keys provided\r\n"), nil
	}

	count := ds.Del(cmd.Params)
	return []byte(":" + strconv.Itoa(count) + "\r\n"), nil
}

// Unknown returns an unknown command error response.
func Unknown(cmd *resp.Command) ([]byte, error) {
	return []byte("-Unknown command " + cmd.Name + "\r\n"), nil
}
