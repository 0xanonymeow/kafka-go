package utils

import (
	"bytes"
	"encoding/gob"
	"io"
	"testing"

	"github.com/0xanonymeow/kafka-go/testings"
)

func TestDeserialize(t *testing.T) {
	type deserializedData struct {
		ID      int
		Name    string
		Balance float64
	}
	type wrongData struct {
		ID   int
		Name string
	}

	var buffer bytes.Buffer
	var emptyBuffer bytes.Buffer

	data := deserializedData{
		ID:      1,
		Name:    "Alice",
		Balance: 95.5,
	}

	subtests := []testings.Subtest{
		{
			Name:         "deserialize",
			ExpectedData: nil,
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				encoder := gob.NewEncoder(&buffer)
				encoder.Encode(data)

				err := Deserialize(&buffer, &deserializedData{})

				return nil, err
			},
		},
		{
			Name:         "deserialize_error",
			ExpectedData: nil,
			ExpectedErr:  io.EOF,
			Test: func() (interface{}, error) {
				err := Deserialize(&emptyBuffer, &deserializedData{})

				return nil, err
			},
		},
		{
			Name:         "deserialize_mismatched",
			ExpectedData: nil,
			ExpectedErr:  io.EOF,
			Test: func() (interface{}, error) {
				err := Deserialize(&buffer, &wrongData{})

				return nil, err
			},
		},
	}

	testings.RunSubtests(t, subtests)
}
