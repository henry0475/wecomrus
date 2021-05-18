package hash

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func GetDestID(dest string) string {
	t := sha1.New()
	io.WriteString(t, dest)
	return fmt.Sprintf("%x", t.Sum(nil))
}
