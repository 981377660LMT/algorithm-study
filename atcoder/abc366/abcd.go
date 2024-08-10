package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

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

type E = int

func (*PreSum3D) e() E        { return 0 }
func (*PreSum3D) op(a, b E) E { return a + b }
func (*PreSum3D) inv(a E) E   { return -a }

type PreSum3D struct {
	xSize, ySize, zSize int
	preSum              [][][]E
}

func NewPreSum3D(mat [][][]E) *PreSum3D {
	res := &PreSum3D{}
	xSize, ySize, zSize := len(mat), len(mat[0]), len(mat[0][0])
	preSum := make([][][]E, xSize+1)
	for x := range preSum {
		preSum[x] = make([][]E, ySize+1)
		for y := range preSum[x] {
			row := make([]E, zSize+1)
			preSum[x][y] = row
		}
	}

	for x := 1; x <= xSize; x++ {
		for y := 1; y <= ySize; y++ {
			for z := 1; z <= zSize; z++ {
				preSum[x][y][z] = mat[x-1][y-1][z-1]
			}
		}
	}

	for x := 1; x <= xSize; x++ {
		for y := 1; y <= ySize; y++ {
			for z := 1; z <= zSize; z++ {
				preSum[x][y][z] = res.op(preSum[x][y][z], preSum[x-1][y][z])
			}
		}
	}

	for x := 1; x <= xSize; x++ {
		for y := 1; y <= ySize; y++ {
			for z := 1; z <= zSize; z++ {
				preSum[x][y][z] = res.op(preSum[x][y][z], preSum[x][y-1][z])
			}
		}
	}

	for x := 1; x <= xSize; x++ {
		for y := 1; y <= ySize; y++ {
			for z := 1; z <= zSize; z++ {
				preSum[x][y][z] = res.op(preSum[x][y][z], preSum[x][y][z-1])
			}
		}
	}

	res.xSize, res.ySize, res.zSize = xSize, ySize, zSize
	res.preSum = preSum
	return res
}

func (ps *PreSum3D) Query(x1, y1, z1, x2, y2, z2 int) E {
	res := ps.preSum[x2+1][y2+1][z2+1]
	res = ps.op(res, ps.inv(ps.preSum[x1][y2+1][z2+1]))
	res = ps.op(res, ps.inv(ps.preSum[x2+1][y1][z2+1]))
	res = ps.op(res, ps.inv(ps.preSum[x2+1][y2+1][z1]))
	res = ps.op(res, ps.preSum[x1][y1][z2+1])
	res = ps.op(res, ps.preSum[x1][y2+1][z1])
	res = ps.op(res, ps.preSum[x2+1][y1][z1])
	res = ps.op(res, ps.inv(ps.preSum[x1][y1][z1]))
	return res
}

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	mat := make([][][]E, n)
	for i := 0; i < n; i++ {
		mat[i] = make([][]E, n)
		for j := 0; j < n; j++ {
			mat[i][j] = make([]E, n)
			for k := 0; k < n; k++ {
				mat[i][j][k] = E(io.NextInt())
			}
		}
	}

	ps := NewPreSum3D(mat)
	q := io.NextInt()
	for i := 0; i < q; i++ {
		lx, rx, ly, ry, lz, rz := io.NextInt(), io.NextInt(), io.NextInt(), io.NextInt(), io.NextInt(), io.NextInt()
		res := ps.Query(lx-1, ly-1, lz-1, rx-1, ry-1, rz-1)
		io.Println(res)
	}

}
