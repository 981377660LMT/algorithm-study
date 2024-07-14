# https://atcoder.jp/contests/abc233/tasks/abc233_g
# 简要题意:给定一张NxN (1<N <50）网格，
# 里面若干个格子有障碍'#'。
# 每次可以选择一个D×D的区域然后消耗D点体力将其中的障碍消去。
# 问将所有障碍消去的最小体力消耗。
# O(n^5)

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n = int(input())
    grid = []
    for _ in range(n):
        grid.append(list(input()))  # '.':空地 '#'：障碍

    def genHash(r1: int, c1: int, r2: int, c2: int) -> int:
        """生成区域的哈希值"""
        return ((r1 * n + c1) * n + r2) * n + c2

    # !TLE
    # def dfs(r1: int, c1: int, r2: int, c2: int) -> int:
    #     """计算区域(r1,c1)到(r2,c2)的最小消耗"""
    #     if r1 == r2 and c1 == c2:
    #         return 0 if grid[r1][c1] == "." else 1
    #     hash_ = genHash(r1, c1, r2, c2)
    #     if ~memo[hash_]:
    #         return memo[hash_]
    #     res = max(r2 - r1 + 1, c2 - c1 + 1)
    #     for r in range(r1, r2):
    #         res = min(res, dfs(r1, c1, r, c2) + dfs(r + 1, c1, r2, c2))
    #     for c in range(c1, c2):
    #         res = min(res, dfs(r1, c1, r2, c) + dfs(r1, c + 1, r2, c2))
    #     memo[hash_] = res
    #     return res

    # memo = [-1] * ((n + 2) * (n + 2) * (n + 2) * (n + 2))
    # res = dfs(0, 0, n - 1, n - 1)
    # print(res)

    dp = [0] * (n * n * n * n)
    for r1 in range(n - 1, -1, -1):
        for c1 in range(n - 1, -1, -1):
            for r2 in range(r1, n):
                for c2 in range(c1, n):
                    hash_ = genHash(r1, c1, r2, c2)
                    if r1 == r2 and c1 == c2:
                        dp[hash_] = 1 if grid[r1][c1] == "#" else 0
                    else:
                        res = max(r2 - r1 + 1, c2 - c1 + 1)
                        for r in range(r1, r2):
                            res = min(
                                res, dp[genHash(r1, c1, r, c2)] + dp[genHash(r + 1, c1, r2, c2)]
                            )
                        for c in range(c1, c2):
                            res = min(
                                res, dp[genHash(r1, c1, r2, c)] + dp[genHash(r1, c + 1, r2, c2)]
                            )
                        dp[hash_] = res

    print(dp[genHash(0, 0, n - 1, n - 1)])

# 21
# .....................
# .....................
# ...#.#...............
# ....#.............#..
# ...#.#...........#.#.
# ..................#..
# .....................
# .....................
# .....................
# ..........#.....#....
# ......#..###.........
# ........#####..#.....
# .......#######.......
# .....#..#####........
# .......#######.......
# ......#########......
# .......#######..#....
# ......#########......
# ..#..###########.....
# .........###.........
# .........###.........
