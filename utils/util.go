package utils

import (
	"fmt"

	"github.com/sesanetwork/go-vassalo/hash"
	"github.com/sesanetwork/go-vassalo/native/idx"
)

// NameOf returns human readable string representation.
func NameOf(p idx.ValidatorID) string {
	if name := hash.GetNodeName(p); len(name) > 0 {
		return name
	}

	return fmt.Sprintf("%d", p)
}
