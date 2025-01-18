package handler_test

import (
	"testing"

	"github.com/Himanshu-Negi8/build-your-own-redis-server/handler"
	"github.com/Himanshu-Negi8/build-your-own-redis-server/types"
)

func TestHandleCommands(t *testing.T) {
	cache := make(map[string]types.CustomValue)

	tests := []struct {
		name     string
		command  []interface{}
		respType types.RESPType
		expected string
	}{
		{
			name:     "PING command",
			command:  []interface{}{"PING"},
			respType: types.RESPTypeSimpleString,
			expected: "+PONG\r\n",
		},
		{
			name:     "ECHO command",
			command:  []interface{}{"ECHO", "Hello, World!"},
			respType: types.RESPTypeArray,
			expected: "$13\r\nHello, World!\r\n",
		},
		{
			name:     "SET command",
			command:  []interface{}{"SET", "mykey", "myvalue"},
			respType: types.RESPTypeArray,
			expected: "+OK\r\n",
		},
		{
			name:     "GET command",
			command:  []interface{}{"GET", "mykey"},
			respType: types.RESPTypeArray,
			expected: "$7\r\nmyvalue\r\n",
		},
		{
			name:     "CONFIG GET dir command",
			command:  []interface{}{"CONFIG", "GET", "dir"},
			respType: types.RESPTypeArray,
			expected: "*2\r\n$3\r\ndir\r\n$15\r\n/tmp/redis-data\r\n",
		},
		{
			name:     "CONFIG GET dbfilename command",
			command:  []interface{}{"CONFIG", "GET", "dbfilename"},
			respType: types.RESPTypeArray,
			expected: "*2\r\n$10\r\ndbfilename\r\n$7\r\nrdbfile\r\n",
		},
		{
			name:     "SAVE command",
			command:  []interface{}{"SAVE"},
			respType: types.RESPTypeArray,
			expected: "+OK\r\n",
		},
		{
			name:     "Invalid command",
			command:  []interface{}{"INVALID"},
			respType: types.RESPTypeArray,
			expected: "-ERR unknown command\r\n",
		},
		{
			name:     "ECHO command with wrong number of arguments",
			command:  []interface{}{"ECHO"},
			respType: types.RESPTypeArray,
			expected: "-ERR wrong number of arguments for 'ECHO' command\r\n",
		},
		{
			name:     "SET command with wrong number of arguments",
			command:  []interface{}{"SET", "mykey"},
			respType: types.RESPTypeArray,
			expected: "-ERR wrong number of arguments for 'SET' command\r\n",
		},
		{
			name:     "SET command with invalid expiration value",
			command:  []interface{}{"SET", "mykey", "myvalue", "PX", "invalid"},
			respType: types.RESPTypeArray,
			expected: "-ERR invalid expiration value\r\n",
		},
		{
			name:     "CONFIG command with unsupported parameter",
			command:  []interface{}{"CONFIG", "GET", "unsupported"},
			respType: types.RESPTypeArray,
			expected: "-ERR unsupported CONFIG parameter\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := handler.HandleCommands(tt.command, tt.respType, cache)
			if string(response) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, response)
			}
		})
	}
}
