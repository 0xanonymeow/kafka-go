package utils

import (
	"bytes"
	"encoding/gob"
)

func Serialize(s interface{}) (bytes.Buffer, error) {
	var b bytes.Buffer

	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(s)

	if err != nil {
		return bytes.Buffer{}, err
	}

	return b, nil
}
