package utils

import (
	"errors"
	"testing"

	"github.com/0xanonymeow/kafka-go/testings"
)

func TestSerialize(t *testing.T) {
	type serializedData struct {
		ID      int
		Name    string
		Balance float64
	}

	data := serializedData{
		ID:      1,
		Name:    "Alice",
		Balance: 95.5,
	}
	wrongData := make(chan int)

	subtests := []testings.Subtest{
		{
			Name:         "serialize",
			ExpectedData: nil,
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				_, err := Serialize(data)

				return nil, err
			},
		},
		{
			Name:         "serialize_error",
			ExpectedData: nil,
			ExpectedErr:  errors.New("gob NewTypeObject can't handle type: chan int"),
			Test: func() (interface{}, error) {
				_, err := Serialize(wrongData)

				return nil, err
			},
		},
	}

	testings.RunSubtests(t, subtests)
}
