# abc379-G - Count Grid 3-coloring (轮廓线dp)
# https://atcoder.jp/contests/abc379/tasks/abc379_g
#
# 给定 h × w 的格子，每个格子写着 ?123中的一种。
# 现将所有的?替换成123中的一个。
# 问有多少替换方法，使得相邻格子的数字不相同。
#
# !H*W<=200
#
# 注意到 H*W<=200，因此最小值<=14.
# dp[row][col][state]表示当前行为row，当前列为col，前col个格子的状态为state的方案数。


from typing import List, Tuple


MOD = 998244353


def countGrid3Coloring(grid: List[str]) -> int:
    h, w = len(grid), len(grid[0])
    if h < w:
        h, w = w, h
        grid = list(zip(*grid))

    newGrid = [[-1] * w for _ in range(h)]  # -1 -> ?, 1/2/3 -> 1/2/3
    for i in range(h):
        for j in range(w):
            if grid[i][j] != "?":
                newGrid[i][j] = int(grid[i][j])

    memo = {}

    def dfs(r: int, c: int, state: Tuple[int, ...]) -> int:
        if r == h:
            return 1
        if c == w:
            return dfs(r + 1, 0, state)

        if (r, c, state) in memo:
            return memo[(r, c, state)]

        cur = newGrid[r][c]
        up = state[0] if r else -2
        left = state[-1] if c else -2
        res = 0
        if cur != -1:
            if cur == up or cur == left:
                res = 0
            else:
                res = dfs(r, c + 1, state[1:] + (cur,))
        else:
            for i in range(1, 4):
                if i == up or i == left:
                    continue
                res += dfs(r, c + 1, state[1:] + (i,))

        res %= MOD
        memo[(r, c, state)] = res
        return res

    res = dfs(0, 0, tuple([-2] * w))
    return res


if __name__ == "__main__":
    H, W = map(int, input().split())
    grid = [input() for _ in range(H)]
    print(countGrid3Coloring(grid))
