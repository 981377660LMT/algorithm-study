package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

// from https://atcoder.jp/users/ccppjsrb
var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	x, y, z := io.NextInt(), io.NextInt(), io.NextInt()
	s := io.Text()

	INF := int(1e18)
	dist := make([][]int, len(s)+1)
	for i := 0; i < len(s)+1; i++ {
		dist[i] = make([]int, 2)
		for j := 0; j < 2; j++ {
			dist[i][j] = INF
		}
	}

	dist[0][0] = 0
	queue := make([][4]int, 0)
	queue = append(queue, [4]int{0, 0, 0, 0}) // dist, caps, pos,preChange
	for len(queue) > 0 {
		curDist, curCaps, curPos, preChange := queue[0][0], queue[0][1], queue[0][2], queue[0][3]
		queue = queue[1:]
		if curPos == len(s) {
			continue
		}
		if dist[curPos][curCaps] < curDist {
			continue
		}

		cand1 := curDist + z
		if preChange == 0 && dist[curPos][1^curCaps] > cand1 {
			dist[curPos][1^curCaps] = cand1
			queue = append(queue, [4]int{cand1, 1 ^ curCaps, curPos, 1})
		}

		upper := 0
		if s[curPos] >= 'A' && s[curPos] <= 'Z' {
			upper = 1
		}
		cand2 := curDist + x
		if upper == curCaps {
			if dist[curPos+1][curCaps] > cand2 {
				dist[curPos+1][curCaps] = cand2
				queue = append(queue, [4]int{cand2, curCaps, curPos + 1, 0})
			}
		}

		cand3 := curDist + y
		if upper != curCaps {
			if dist[curPos+1][curCaps] > cand3 {
				dist[curPos+1][curCaps] = cand3
				queue = append(queue, [4]int{cand3, curCaps, curPos + 1, 0})
			}
		}
	}

	res := INF
	for i := 0; i < 2; i++ {
		if res > dist[len(s)][i] {
			res = dist[len(s)][i]
		}
	}
	io.Println(res)

}
