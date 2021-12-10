package utils

import (
	"io"
	"log"
	"os"
)

func Copy(from, to string) {
	src, err := os.Open(from)
	if err != nil {
		log.Panic(err)
	}
	defer src.Close()
	dst, err := os.Create(to)
	if err != nil {
		log.Panic(err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		log.Panic(err)
	}
}

// Exist uses err returned from os.Stat to determine if a file/folder exists
func Exist(err error) bool {
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			log.Panic(err)
		}
	}
	return true
}
