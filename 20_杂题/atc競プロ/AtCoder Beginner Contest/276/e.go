// import sys

// sys.setrecursionlimit(int(1e9))
// input = lambda: sys.stdin.readline().rstrip("\r\n")
// MOD = 998244353
// INF = int(4e18)

// # 縦 H 行、横 W 列のマス目があり、上から i(1≤i≤H) 行目、左から j(1≤j≤W) 列目のマスを (i,j) と表します。

// # 各マスは「始点」「道」「障害物」のいずれかです。
// # マス (i,j) の状態は文字 C
// # i,j
// # ​
// #   で表され、C
// # i,j
// # ​
// #  = S なら始点、C
// # i,j
// # ​
// #  = . なら道、C
// # i,j
// # ​
// #  = # なら障害物です。始点のマスはただ一つ存在します。

// # 始点のマスを出発し、上下または左右に隣接するマスに移動することを繰り返して、障害物のマスを通らずに始点のマスへ戻ってくるような長さ 4 以上の経路であって、最初と最後を除き同じマスを通らないようなものが存在するか判定してください。
// # より厳密には、以下の条件を満たす整数 n およびマスの列 (x
// # 0
// # ​
// #  ,y
// # 0
// # ​
// #  ),(x
// # 1
// # ​
// #  ,y
// # 1
// # ​
// #  ),…,(x
// # n
// # ​
// #  ,y
// # n
// # ​
// #  ) が存在するか判定してください。

// # n≥4
// # C
// # x
// # 0
// # ​
// #  ,y
// # 0
// # ​

// # ​
// #  =C
// # x
// # n
// # ​
// #  ,y
// # n
// # ​

// # ​
// #  = S
// # 1≤i≤n−1 ならば C
// # x
// # i
// # ​
// #  ,y
// # i
// # ​

// # ​
// #  = .
// # 1≤i<j≤n−1 ならば (x
// # i
// # ​
// #  ,y
// # i
// # ​
// #  )
// # 
// # =(x
// # j
// # ​
// #  ,y
// # j
// # ​
// #  )
// # 0≤i≤n−1 ならばマス (x
// # i
// # ​
// #  ,y
// # i
// # ​
// #  ) とマス (x
// # i+1
// # ​
// #  ,y
// # i+1
// # ​
// #  ) は上下または左右に隣接する

// # 4 4
// # ....
// # #.#.
// # .S..
// # .##.
// # Yes

// DIR4 = [(1, 0), (-1, 0), (0, 1), (0, -1)]
// !是否在大小>=4的环上 (不能并查集/拓扑排序 需要dfs)
// if __name__ == "__main__":
//     ROW, COL = map(int, input().split())
//     grid = [list(input()) for _ in range(ROW)]

//     sr, sc = -1, -1
//     for r in range(ROW):
//         for c in range(COL):
//             if grid[r][c] == "S":
//                 sr, sc = r, c
//                 break
//         if sr != -1:
//             break

//     def dfs(cur: int, depth: int) -> None:
//         global res
//         if visited[cur]:
//             if cur == start and depth >= 4:
//                 res = True
//             return
//         visited[cur] = True
//         for dr, dc in DIR4:
//             nr, nc = cur // COL + dr, cur % COL + dc
//             if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] != "#":
//                 dfs(nr * COL + nc, depth + 1)

//     res = False
//     start = sr * COL + sc
//     visited = [False] * (ROW * COL)
//     dfs(start, 0)
//     print("Yes" if res else "No")
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ROW, COL int
	fmt.Fscan(in, &ROW, &COL)
	grid := make([]string, ROW)
	for i := 0; i < ROW; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = s
	}

	var sr, sc int
	for r := 0; r < ROW; r++ {
		for c := 0; c < COL; c++ {
			if grid[r][c] == 'S' {
				sr, sc = r, c
				break
			}
		}
	}

	dir4 := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	var res bool
	start := sr*COL + sc
	visited := make([]bool, ROW*COL)
	var dfs func(cur, depth int)

	dfs = func(cur, depth int) {
		if visited[cur] {
			if cur == start && depth >= 4 {
				res = true
			}
			return
		}

		visited[cur] = true
		for _, d := range dir4 {
			nr, nc := cur/COL+d[0], cur%COL+d[1]
			if 0 <= nr && nr < ROW && 0 <= nc && nc < COL && grid[nr][nc] != '#' {
				dfs(nr*COL+nc, depth+1)
			}
		}
	}

	dfs(start, 0)
	if res {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}

}
