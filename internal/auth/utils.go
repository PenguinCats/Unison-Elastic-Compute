package auth

import (
	"crypto/rand"
	"github.com/gofrs/uuid"
	"math/big"
)

func GenerateRandomUUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}

func GenerateRandomInt() int64 {
	v, _ := rand.Int(rand.Reader, big.NewInt(19981125))
	return v.Int64()
}
