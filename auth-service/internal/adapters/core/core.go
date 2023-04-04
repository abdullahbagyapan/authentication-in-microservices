package core

import (
	"crypto/sha512"
	"fmt"
)

type Adapter struct {
}

func NewAdapter() *Adapter {
	return new(Adapter)
}

func (A Adapter) Hash(value string) string {

	byteValue := []byte(value)

	hashedValue := sha512.Sum384(byteValue)

	str := fmt.Sprintf("%x", hashedValue)

	return str

}
