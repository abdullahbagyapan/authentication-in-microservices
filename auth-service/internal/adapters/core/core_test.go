package core

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {

	coreAdapter := NewAdapter()

	hashedValue := coreAdapter.Hash("12asdxcz")

	fmt.Println(hashedValue)

}
