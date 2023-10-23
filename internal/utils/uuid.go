package utils

import (
	"encoding/binary"

	"github.com/google/uuid"
)

func Int64ToUuid(prefix, id uint64) (string, error) {
	a, b := make([]byte, 8), make([]byte, 8)
	binary.LittleEndian.PutUint64(a, prefix)
	binary.LittleEndian.PutUint64(b, id)

	u, err := uuid.FromBytes(append(a, b...))
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
