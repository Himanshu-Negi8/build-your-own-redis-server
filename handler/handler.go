package handler

import (
	"fmt"
	"github.com/Himanshu-Negi8/build-your-own-redis-server/types"
	"os"
	"strconv"
	"time"
)

const (
	dir        = "/tmp/redis-data"
	dbfilename = "rdbfile"
)

func HandleCommands(commandTokens interface{}, respType types.RESPType, cache map[string]types.CustomValue) []byte {
	// Process the command based on its type
	switch respType {
	case types.RESPTypeArray:
		arr := commandTokens.([]interface{})
		if len(arr) >= 1 {
			switch arr[0].(string) {
			case "PING":
				return []byte("+PONG\r\n")
			case "ECHO":
				return echoCommand(arr)
			case "SET":
				return setCommand(arr, cache)
			case "GET":
				return getCommand(arr, cache)
			case "CONFIG":
				return configCommand(arr)
			case "SAVE":
				return saveCommand(arr, cache)
			}

		} else {
			return []byte("-ERR unknown command\r\n")
		}
	case types.RESPTypeSimpleString:
		return []byte("+OK\r\n")
	default:
		return []byte("-ERR unsupported RESP type\r\n")
	}

	return []byte("-ERR unsupported RESP type\r\n")
}

func echoCommand(arr []interface{}) []byte {
	// validate length is exactly 2
	// ECHO <message>
	if len(arr) != 2 {
		return []byte("-ERR wrong number of arguments for 'ECHO' command\r\n")
	}

	message := arr[1].(string)
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(message), message))
}

func setCommand(arr []interface{}, cache map[string]types.CustomValue) []byte {
	// validate length is exactly 3 or 5
	// SET <key> <value>
	// SET <key> <value> <px> <expiration>
	// In RESP we will receive the key as *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n which is an array of 3 elements
	key := arr[1].(string)
	value := arr[2].(string)
	if len(arr) == 5 {
		expiration, err := strconv.ParseInt(arr[4].(string), 10, 64)
		if err != nil {
			return []byte("-ERR invalid expiration value\r\n")

		}
		cache[key] = types.CustomValue{Value: value, ValueExpiration: time.Now().UnixMilli() + expiration}
	} else if len(arr) == 3 {
		// This is without the expiration time.
		cache[key] = types.CustomValue{Value: value, ValueExpiration: -1}
	} else {
		return []byte("-ERR wrong number of arguments for 'SET' command\r\n")
	}

	return []byte(fmt.Sprintf("+OK\r\n"))
}

func getCommand(arr []interface{}, cache map[string]types.CustomValue) []byte {
	// validate length is exactly 2
	// GET <key>
	if len(arr) != 2 {
		return []byte("-ERR wrong number of arguments for 'GET' command\r\n")
	}

	key := arr[1].(string)
	val, ok := cache[key]

	if ok && (val.ValueExpiration == -1 || val.ValueExpiration > time.Now().UnixMilli()) {
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val.Value), val.Value))
	} else {
		return []byte("$-1\r\n")
	}
}

func configCommand(arr []interface{}) []byte {
	if len(arr) != 3 {
		return []byte("-ERR wrong number of arguments for 'CONFIG' command\r\n")
	}

	if arr[1].(string) == "GET" && arr[2].(string) == "dir" {
		return []byte(fmt.Sprintf("*2\r\n$3\r\ndir\r\n$%d\r\n%s\r\n", len(dir), dir))
	} else if arr[1].(string) == "GET" && arr[2].(string) == "dbfilename" {
		return []byte(fmt.Sprintf("*2\r\n$10\r\ndbfilename\r\n$%d\r\n%s\r\n", len(dbfilename), dbfilename))
	}
	return []byte("-ERR unsupported CONFIG parameter\r\n")
}

func saveCommand(arr []interface{}, cache map[string]types.CustomValue) []byte {
	// validate length is exactly 1
	// SAVE
	if len(arr) != 1 {
		return []byte("-ERR wrong number of arguments for 'SAVE' command\r\n")
	}

	fmt.Println("Saving data to file")
	// Create the directory if it doesn't exist
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return []byte("-ERR failed to save data to file\r\n")
	}

	// Create the file
	filePath := fmt.Sprintf("%s/%s", dir, dbfilename)
	file, err := os.Create(filePath)
	if err != nil {
		return []byte("-ERR failed to save data to file\r\n")
	}

	file.WriteString("REDIS0009")
	defer file.Close()
	//
	// Ensure data is flushed to disk
	if err := file.Sync(); err != nil {
		fmt.Println("Error syncing file to disk:", err)
		return []byte("-ERR failed to save data to file\r\n")
	}

	fmt.Println("File saved successfully at", filePath)
	return []byte("+OK\r\n")
}
