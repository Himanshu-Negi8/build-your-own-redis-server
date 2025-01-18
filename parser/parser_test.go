package parser_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/Himanshu-Negi8/build-your-own-redis-server/parser"
	"github.com/Himanshu-Negi8/build-your-own-redis-server/types"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedResult interface{}
		expectedType   types.RESPType
		expectedError  error
	}{
		{
			name:           "Simple String",
			input:          "+OK\r\n",
			expectedResult: "OK",
			expectedType:   types.RESPTypeSimpleString,
			expectedError:  nil,
		},
		{
			name:           "Integer",
			input:          ":1000\r\n",
			expectedResult: 1000,
			expectedType:   types.RESPTypeInteger,
			expectedError:  nil,
		},
		{
			name:           "Bulk String",
			input:          "$6\r\nfoobar\r\n",
			expectedResult: "foobar",
			expectedType:   types.RESPTypeBulkString,
			expectedError:  nil,
		},
		{
			name:           "Null Bulk String",
			input:          "$-1\r\n",
			expectedResult: nil,
			expectedType:   types.RESPTypeBulkString,
			expectedError:  nil,
		},
		//{
		//	name:           "Array",
		//	input:          "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
		//	expectedResult: []interface{}{"foo", "bar"},
		//	expectedType:   types.RESPTypeArray,
		//	expectedError:  nil,
		//},
		{
			name:           "Null Array",
			input:          "*-1\r\n",
			expectedResult: nil,
			expectedType:   types.RESPTypeArray,
			expectedError:  nil,
		},
		{
			name:           "Unknown Prefix",
			input:          "!unknown\r\n",
			expectedResult: nil,
			expectedType:   "",
			expectedError:  errors.New("unknown prefix: !"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.input))
			result, respType, err := parser.Parse(reader)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result != tt.expectedResult {
					t.Errorf("expected result %v, got %v", tt.expectedResult, result)
				}
				if respType != tt.expectedType {
					t.Errorf("expected respType %v, got %v", tt.expectedType, respType)
				}
			}
		})
	}
}
