package utils

import (
	"bytes"
	"encoding/gob"
)

func Deserialize(b *bytes.Buffer, s interface{}) error {
	decoder := gob.NewDecoder(b)
	err := decoder.Decode(s)

	if err != nil {
		return err
	}

	return nil
}
