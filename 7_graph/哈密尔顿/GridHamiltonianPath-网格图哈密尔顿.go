// 始点を固定して、ハミルトンパスを作る。偶数ならサイクルにする。
// 高次元版：https://codeforces.com/contest/863/submission/194294053

// 给定起点，求网格中的哈密尔顿路径/哈密尔顿回路

package main

import "fmt"

func main() {
	fmt.Println(GridHamiltonianPath(2, 2, 0, 0))
	fmt.Println(GridHamiltonianPath(3, 3, 1, 0))
}

func GridHamiltonianPath(H, W, sx, sy int) [][2]int {
	if H == 1 {
		path := [][2]int{}
		if sy == 0 {
			for y := 0; y < W; y++ {
				path = append(path, [2]int{0, y})
			}
			return path
		}
		if sy == W-1 {
			for y := W - 1; y >= 0; y-- {
				path = append(path, [2]int{0, y})
			}
			return path
		}
		return nil
	}

	if W == 1 {
		path := [][2]int{}
		if sx == 0 {
			for x := 0; x < H; x++ {
				path = append(path, [2]int{x, 0})
			}
			return path
		}
		if sx == H-1 {
			for x := H - 1; x >= 0; x-- {
				path = append(path, [2]int{x, 0})
			}
			return path
		}
		return nil
	}

	if H%2 == 0 && W%2 == 1 {
		path := GridHamiltonianPath(W, H, sy, sx)
		for _, p := range path {
			p[0], p[1] = p[1], p[0]
		}
		return path
	}

	if W%2 == 0 {
		path := make([][2]int, 0, H*W)
		for j := 0; j < W; j++ {
			path = append(path, [2]int{0, j})
		}
		for j := W - 1; j >= 0; j-- {
			if j%2 == 1 {
				for i := 1; i < H; i++ {
					path = append(path, [2]int{i, j})
				}
			}
			if j%2 == 0 {
				for i := H - 1; i >= 1; i-- {
					path = append(path, [2]int{i, j})
				}
			}
		}

		idx := -1
		for i, p := range path {
			if p[0] == sx && p[1] == sy {
				idx = i
				break
			}
		}
		if idx == -1 {
			return nil
		}
		path = append(path[idx:], path[:idx]...)
		return path
	}

	if (sx+sy)%2 == 1 {
		return nil
	}

	if sx%2 == 1 {
		path := make([][2]int, 0, H*W)
		for i := sx; i >= 0; i-- {
			if i%2 == 1 {
				for j := sy; j >= 0; j-- {
					path = append(path, [2]int{i, j})
				}
			}
			if i%2 == 0 {
				for j := 0; j <= sy; j++ {
					path = append(path, [2]int{i, j})
				}
			}
		}
		for j := sy + 1; j < W; j++ {
			if j%2 == 0 {
				for i := 0; i <= sx; i++ {
					path = append(path, [2]int{i, j})
				}
			}
			if j%2 == 1 {
				for i := sx; i >= 0; i-- {
					path = append(path, [2]int{i, j})
				}
			}
		}
		for j := W - 1; j >= 0; j-- {
			if j%2 == 0 {
				for i := sx + 1; i < H; i++ {
					path = append(path, [2]int{i, j})
				}
			}
			if j%2 == 1 {
				for i := H - 1; i > sx; i-- {
					path = append(path, [2]int{i, j})
				}
			}
		}
		return path
	}

	if sx == H-1 {
		path := GridHamiltonianPath(H, W, 0, sy)
		for _, p := range path {
			p[0] = H - 1 - p[0]
		}
		return path
	}

	if sy == W-1 {
		path := GridHamiltonianPath(H, W, sx, 0)
		for _, p := range path {
			p[1] = W - 1 - p[1]
		}
		return path
	}

	if sx != 0 && sy == 0 {
		path := GridHamiltonianPath(W, H, sy, sx)
		for _, p := range path {
			p[0], p[1] = p[1], p[0]
		}
		return path
	}

	path := make([][2]int, 0, H*W)
	if sx == 0 {
		for j := sy; j >= 0; j-- {
			if j%2 == 0 {
				for i := 0; i < H-1; i++ {
					path = append(path, [2]int{i, j})
				}
			}
			if j%2 == 1 {
				for i := H - 2; i >= 0; i-- {
					path = append(path, [2]int{i, j})
				}
			}
		}
		for j := 0; j < W; j++ {
			path = append(path, [2]int{H - 1, j})
		}
		for j := W - 1; j > sy; j-- {
			if j%2 == 0 {
				for i := H - 2; i >= 0; i-- {
					path = append(path, [2]int{i, j})
				}
			}
			if j%2 == 1 {
				for i := 0; i < H-1; i++ {
					path = append(path, [2]int{i, j})
				}
			}
		}
		return path
	}

	for j := sy; j >= 0; j-- {
		path = append(path, [2]int{sx, j})
	}

	for j := 0; j < sy+1; j++ {
		if j%2 == 0 {
			for i := sx - 1; i >= 0; i-- {
				path = append(path, [2]int{i, j})
			}
		}
		if j%2 == 1 {
			for i := 0; i < sx; i++ {
				path = append(path, [2]int{i, j})
			}
		}
	}

	for i := 0; i < sx+1; i++ {
		if i%2 == 0 {
			for j := sy + 1; j < W; j++ {
				path = append(path, [2]int{i, j})
			}
		}
		if i%2 == 1 {
			for j := W - 1; j > sy; j-- {
				path = append(path, [2]int{i, j})
			}
		}
	}

	for j := W - 1; j >= 0; j-- {
		if j%2 == 0 {
			for i := sx + 1; i < H; i++ {
				path = append(path, [2]int{i, j})
			}
		}
		if j%2 == 1 {
			for i := H - 1; i > sx; i-- {
				path = append(path, [2]int{i, j})
			}
		}
	}
	return path
}
