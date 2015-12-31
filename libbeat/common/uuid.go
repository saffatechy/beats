package common

import (
	"crypto/rand"
	"fmt"
)

// GenerateUUID is used to generate a random UUID
// This was copied from https://github.com/hashicorp/uuid/blob/master/uuid.go
func UUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16])
}
