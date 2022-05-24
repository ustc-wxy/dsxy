package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
)

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	if _, e := io.Copy(h, r); e != nil {
		log.Fatal(e)
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
