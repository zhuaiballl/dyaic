package utils

import (
	"bytes"
	"crypto/md5"
	"io"
	"log"
	"os"
)

func Md5File(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	m := md5.New()
	_, err = io.Copy(m, file)
	if err != nil {
		log.Panic(err)
	}

	return m.Sum(nil)
}

func SameFile(a, b string) bool {
	md5A := Md5File(a)
	md5B := Md5File(b)
	return bytes.Compare(md5A, md5B) == 0
}
