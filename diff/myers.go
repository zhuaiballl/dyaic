package diff

import (
	"bufio"
	"dyaic/config"
	"fmt"
	"io"
	"log"
	"os"
)

type operation uint

type Change struct {
	LineNumber int
	Op         operation
	Content    string
}

type Changes struct {
	Item []Change
}

const (
	INSERT operation = 1
	DELETE operation = 2
	MOVE   operation = 3
)

func (op operation) String() string {
	switch op {
	case INSERT:
		return "INS"
	case DELETE:
		return "DEL"
	case MOVE:
		return "MOV"
	default:
		return "UNKNOWN"
	}
}

var colors = map[operation]string{
	INSERT: "\033[32m",
	DELETE: "\033[31m",
	MOVE:   "\033[39m",
}

func ShowDiff(old, new string) {
	src, err := getFileLines(old)
	if err != nil {
		log.Panic(err)
	}

	dst, err := getFileLines(new)
	if err != nil {
		log.Panic(err)
	}

	GenerateDiff(src, dst)
}

func GenerateChanges(old, new string) (chs Changes) {
	src, err := getFileLines(old)
	if err != nil {
		log.Panic(err)
	}
	dst, err := getFileLines(new)
	if err != nil {
		log.Panic(err)
	}
	script := shortestEditScript(src, dst)
	srcIndex, dstIndex := 0, 0

	for _, op := range script {
		switch op {
		case INSERT:
			chs.Item = append(chs.Item, Change{
				LineNumber: srcIndex,
				Op:         INSERT,
				Content:    dst[dstIndex],
			})
			dstIndex++
		case MOVE:
			srcIndex++
			dstIndex++
		case DELETE:
			chs.Item = append(chs.Item, Change{
				LineNumber: srcIndex,
				Op:         DELETE,
				Content:    "",
			})
			srcIndex++
		}
	}
	return
}

func GenerateDiff(src, dst []string) {
	script := shortestEditScript(src, dst)

	srcIndex, dstIndex := 0, 0

	for _, op := range script {
		switch op {
		case INSERT:
			fmt.Println(colors[op] + "+" + dst[dstIndex])
			dstIndex += 1

		case MOVE:
			fmt.Println(colors[op] + " " + src[srcIndex])
			srcIndex += 1
			dstIndex += 1

		case DELETE:
			fmt.Println(colors[op] + "-" + src[srcIndex])
			srcIndex += 1
		}
	}
	fmt.Println(colors[MOVE])
}

func SaveDiff(old, new string) {
	src, err := getFileLines(old)
	if err != nil {
		log.Panic(err)
	}

	dst, err := getFileLines(new)
	if err != nil {
		log.Panic(err)
	}
	script := shortestEditScript(src, dst)
	f, err := os.Create(config.RepoLocation + "/diff")
	defer f.Close()
	if err != nil {
		log.Panic(err)
	}
	srcIndex, dstIndex := 0, 0

	for _, op := range script {
		switch op {
		case INSERT:
			io.WriteString(f, "+"+dst[dstIndex])
			dstIndex += 1

		case MOVE:
			//fmt.Println(colors[op] + " " + src[srcIndex])
			srcIndex += 1
			dstIndex += 1

		case DELETE:
			io.WriteString(f, "-"+src[srcIndex])
			srcIndex += 1
		}
	}
}

func shortestEditScript(src, dst []string) []operation {
	n := len(src)
	m := len(dst)
	max := n + m
	var trace []map[int]int
	var x, y int

loop:
	for d := 0; d <= max; d++ {
		v := make(map[int]int, d+2)
		trace = append(trace, v)

		if d == 0 {
			t := 0
			for len(src) > t && len(dst) > t && src[t] == dst[t] {
				t++
			}
			v[0] = t
			if t == len(src) && t == len(dst) {
				break loop
			}
			continue
		}

		lastV := trace[d-1]

		for k := -d; k <= d; k += 2 {
			if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
				x = lastV[k+1]
			} else {
				x = lastV[k-1] + 1
			}

			y = x - k

			for x < n && y < m && src[x] == dst[y] {
				x, y = x+1, y+1
			}

			v[k] = x

			if x == n && y == m {
				break loop
			}
		}
	}

	var script []operation

	x = n
	y = m
	var k, prevK, prevX, prevY int

	for d := len(trace) - 1; d > 0; d-- {
		k = x - y
		lastV := trace[d-1]

		if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX = lastV[prevK]
		prevY = prevX - prevK

		for x > prevX && y > prevY {
			script = append(script, MOVE)
			x -= 1
			y -= 1
		}

		if x == prevX {
			script = append(script, INSERT)
		} else {
			script = append(script, DELETE)
		}

		x, y = prevX, prevY
	}

	if trace[0][0] != 0 {
		for i := 0; i < trace[0][0]; i++ {
			script = append(script, MOVE)
		}
	}

	return reverse(script)
}

func printTrace(trace []map[int]int) {
	for d := 0; d < len(trace); d++ {
		fmt.Printf("d = %d:\n", d)
		v := trace[d]
		for k := -d; k <= d; k += 2 {
			x := v[k]
			y := x - k
			fmt.Printf("  k = %2d: (%d, %d)\n", k, x, y)
		}
	}
}

func reverse(s []operation) []operation {
	result := make([]operation, len(s))

	for i, v := range s {
		result[len(s)-1-i] = v
	}

	return result
}

func getFileLines(p string) ([]string, error) {
	f, err := os.Open(p)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
