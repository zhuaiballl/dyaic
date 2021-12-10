package diff

import (
	"log"
	"os"
)

func Recover(dst string, chs *Changes) error {
	lines, err := getFileLines(dst)
	if err != nil {
		log.Panic(err)
	}
	index := 0
	f, err := os.Create(dst)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	for _, ch := range chs.Item {
		for index < ch.LineNumber {
			_, err = f.WriteString(lines[index] + "\n")
			if err != nil {
				log.Panic(err)
			}
			index++
		}
		switch ch.Op {
		case INSERT:
			_, err = f.WriteString(ch.Content + "\n")
			if err != nil {
				log.Panic(err)
			}
		case DELETE:
			index++
		}
	}
	for index < len(lines) {
		_, err = f.WriteString(lines[index] + "\n")
		if err != nil {
			log.Panic(err)
		}
		index++
	}
	return nil
}
