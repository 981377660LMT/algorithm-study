package main

import "fmt"

func main() {
	G := NewFillGridWithL()
	res := G.Solve(5, 6)
	fmt.Println(res)
	color := G.Color(res)
	fmt.Println(color)
}

var DIR4 = [][]int32{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

// 蒙德里安的梦想：用L形块填充格子.
// 多米诺平铺.
type FillGridWithL struct {
	res       [][]int32
	mat       [][]int32
	color     [][]int32
	indicator int32
}

func NewFillGridWithL() *FillGridWithL {
	return &FillGridWithL{}
}

func (f *FillGridWithL) Solve(n, m int32) [][]int32 {
	f.indicator = 0
	flip := false
	if n%3 == 0 && m%3 != 0 {
		n, m = m, n
		flip = !flip
	}
	if n%6 == 0 {
		n, m = m, n
		flip = !flip
	}
	if n <= 1 || m <= 1 || m%3 != 0 {
		return nil
	}
	f.res = make([][]int32, n)
	for i := int32(0); i < n; i++ {
		f.res[i] = make([]int32, m)
	}
	if n%2 == 0 {
		for i := int32(0); i < n; i += 2 {
			for j := int32(0); j < m; j += 3 {
				f.paint4(i, j)
				f.paint2(i+1, j+2)
			}
		}
	} else if m%2 == 0 {
		for i := int32(0); i < n-3; i += 2 {
			for j := int32(0); j < m; j += 3 {
				f.paint4(i, j)
				f.paint2(i+1, j+2)
			}
		}
		for i := int32(0); i < m; i += 2 {
			f.paint4(n-3, i)
			f.paint2(n-1, i+1)
		}
	} else if m >= 9 && n >= 5 {
		f.paint4(0, 0)
		f.paint4(2, 0)
		f.paint2(4, 1)
		f.paint3(1, 2)
		f.paint3(0, 3)
		f.paint1(4, 2)
		for i := int32(3); i+2 <= m-2; i += 2 {
			f.paint4(2, i)
			f.paint1(4, i+1)
		}
		f.paint3(3, m-1)
		f.paint2(2, m-1)
		f.paint4(0, m-2)
		for i := int32(4); i < m-2; i += 3 {
			f.paint1(1, i)
			f.paint3(0, i+2)
		}
		for i := int32(5); i < n; i += 2 {
			for j := int32(0); j < m; j += 3 {
				f.paint4(i, j)
				f.paint2(i+1, j+2)
			}
		}
	} else {
		return nil
	}
	if flip {
		f.res = f.transpose(f.res)
	}
	return f.res
}

func (f *FillGridWithL) Color(mat [][]int32) [][]int32 {
	n, m := int32(len(mat)), int32(len(mat[0]))
	f.mat = mat
	f.color = make([][]int32, n)
	for i := int32(0); i < n; i++ {
		f.color[i] = make([]int32, m)
		for j := int32(0); j < m; j++ {
			f.color[i][j] = -1
		}
	}
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < m; j++ {
			if f.color[i][j] == -1 {
				bits := f.findColor(i, j, mat[i][j], -1, -1)
				c := int32(0)
				for bits>>c&1 == 1 {
					c++
				}
				f.setColor(i, j, mat[i][j], -1, -1, c)
			}
		}
	}
	return f.color
}

func (f *FillGridWithL) paint1(i, j int32) {
	f.indicator++
	f.res[i][j] = f.indicator
	f.res[i-1][j] = f.indicator
	f.res[i][j+1] = f.indicator
}

func (f *FillGridWithL) paint2(i, j int32) {
	f.indicator++
	f.res[i][j] = f.indicator
	f.res[i-1][j] = f.indicator
	f.res[i][j-1] = f.indicator
}

func (f *FillGridWithL) paint3(i, j int32) {
	f.indicator++
	f.res[i][j] = f.indicator
	f.res[i+1][j] = f.indicator
	f.res[i][j-1] = f.indicator
}

func (f *FillGridWithL) paint4(i, j int32) {
	f.indicator++
	f.res[i][j] = f.indicator
	f.res[i+1][j] = f.indicator
	f.res[i][j+1] = f.indicator
}

func (f *FillGridWithL) transpose(mat [][]int32) [][]int32 {
	n, m := int32(len(mat)), int32(len(mat[0]))
	res := make([][]int32, m)
	for i := int32(0); i < m; i++ {
		res[i] = make([]int32, n)
	}
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < m; j++ {
			res[j][i] = mat[i][j]
		}
	}
	return res
}

func (f *FillGridWithL) valid(i, j int32) bool {
	return i >= 0 && j >= 0 && i < int32(len(f.mat)) && j < int32(len(f.mat[0]))
}

func (f *FillGridWithL) setColor(i, j, v, fi, fj, c int32) {
	if v != f.mat[i][j] {
		return
	}
	f.color[i][j] = c
	for _, d := range DIR4 {
		x, y := d[0]+i, d[1]+j
		if !f.valid(x, y) {
			continue
		}
		if x == fi && y == fj {
			continue
		}
		f.setColor(x, y, v, i, j, c)
	}
}

func (f *FillGridWithL) findColor(i, j, v, fi, fj int32) int32 {
	if v != f.mat[i][j] {
		if f.color[i][j] == -1 {
			return 0
		}
		return 1 << f.color[i][j]
	}
	res := int32(0)
	for _, d := range DIR4 {
		x, y := d[0]+i, d[1]+j
		if !f.valid(x, y) {
			continue
		}
		if x == fi && y == fj {
			continue
		}
		res |= f.findColor(x, y, v, i, j)
	}
	return res
}
