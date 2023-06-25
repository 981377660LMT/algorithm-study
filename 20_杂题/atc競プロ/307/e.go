package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"regexp"
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

	io.NextInt()
	s := io.Text()

	res := regexp.MustCompile(`(?U)\((.*)\)`).ReplaceAllString(s, "")
	io.Println(res)

}
