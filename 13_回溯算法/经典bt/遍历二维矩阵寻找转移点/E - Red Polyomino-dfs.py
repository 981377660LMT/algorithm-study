"""
dfs+状态压缩
"""

# 红色骨牌(七巧板)

# 给你 一个 n*n的矩阵,
# 其中各单元属性要么是白要么是黑,
# 现在需要你计算出有多少种方案可以将`白块绘制成红块`,
# !使红块区域连通且大小为k。
# !n,k<=8
# !"#"表示黑色,"."表示白色


# !直接dfs+visited 存储看过的状态来搜索
# !每次dfs遍历整个矩阵,寻找下个红色块,判断和之前的红色块是否连通

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


DIR4 = [(1, 0), (-1, 0), (0, 1), (0, -1)]


if __name__ == "__main__":

    def dfs(count: int, state: int) -> None:
        global res
        if count == k:
            res += 1
            return

        for r in range(n):
            for c in range(n):
                if grid[r][c] == "#" or state & (1 << (r * n + c)):
                    continue
                ok = False
                for dr, dc in DIR4:
                    nr, nc = r + dr, c + dc
                    if 0 <= nr < n and 0 <= nc < n and state & (1 << (nr * n + nc)):
                        ok = True
                        break
                if ok:
                    nextState = state | (1 << (r * n + c))
                    if nextState not in visited:
                        visited.add(nextState)
                        dfs(count + 1, nextState)

    n = int(input())
    k = int(input())
    grid = [list(input()) for _ in range(n)]

    visited = set()
    res = 0
    for r in range(n):
        for c in range(n):
            if grid[r][c] == "#":
                continue
            visited.add(1 << (r * n + c))
            dfs(1, 1 << (r * n + c))
    print(res)
