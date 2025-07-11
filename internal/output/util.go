package output

import (
	"crypto/sha256"
	"encoding/json"
)

func HashStruct(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(data)
	return hash[:], nil
}
