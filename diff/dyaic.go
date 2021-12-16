package diff

import (
	"dyaic/config"
	"io/ioutil"
	"log"
)

func shortestEditDist(src, dst []byte) []operation {
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

func SaveDyaicDiff(old, new string) {
	src, err := ioutil.ReadFile(old)
	if err != nil {
		log.Panic(err)
	}

	dst, err := ioutil.ReadFile(new)
	if err != nil {
		log.Panic(err)
	}
	script := shortestEditDist(src, dst)
	//f, err := os.Create(config.RepoLocation + "/diff")
	//defer f.Close()
	//if err != nil {
	//	log.Panic(err)
	//}
	srcIndex, dstIndex := 0, 0
	lstOp := operation(0)
	var tmp []byte
	for _, op := range script {
		if op != lstOp {
			if op == INSERT {
				tmp = append(tmp, byte('+'))
			}
			if op == DELETE {
				tmp = append(tmp, byte('-'))
			}
			lstOp = op
		}
		switch op {
		case INSERT:
			//io.WriteString(f, "+"+dst[dstIndex])
			tmp = append(tmp, dst[dstIndex])
			dstIndex += 1

		case MOVE:
			//fmt.Println(colors[op] + " " + src[srcIndex])
			srcIndex += 1
			dstIndex += 1

		case DELETE:
			//io.WriteString(f, "-"+src[srcIndex])
			tmp = append(tmp, src[srcIndex])
			srcIndex += 1
		}
	}
	err = ioutil.WriteFile(config.RepoLocation+"/diff", tmp, 0644)
	if err != nil {
		log.Panic(err)
	}
}
