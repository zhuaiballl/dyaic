package diff

import (
	"crypto/md5"
	"fmt"
	"github.com/jackpal/bencode-go"
	"io"
	"log"
	"os"
	"time"
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

func Md5FileTest(filepath string) {
	beginTime := time.Now()
	file, err := os.Open(filepath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		log.Panic(err)
	}

	m := md5.New()
	_, err = io.Copy(m, file)
	if err != nil {
		log.Panic(err)
	}
	res := m.Sum(nil)
	endTime := time.Now()
	res = nil
	fmt.Println(fileInfo.Size(), endTime.Sub(beginTime).Nanoseconds(), res)
}

func Ben(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	err = bencode.Marshal(os.Stdout, file)
	if err != nil {
		log.Panic(err)
	}
	return nil
}
