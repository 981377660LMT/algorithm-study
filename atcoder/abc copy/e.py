import sys
from collections import defaultdict

sys.setrecursionlimit(int(1e7))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def solve():
    H, W = map(int, input().split())
    grid = [list(input()) for _ in range(H)]

    # 转换为一维
    flat_grid = []
    for row in grid:
        for cell in row:
            if cell == "?":
                flat_grid.append(-1)
            else:
                flat_grid.append(int(cell))

    total = H * W

    # 预处理相邻位置
    adj = [[] for _ in range(total)]
    for i in range(H):
        for j in range(W):
            pos = i * W + j
            directions = [(-1, 0), (1, 0), (0, -1), (0, 1)]
            for di, dj in directions:
                ni, nj = i + di, j + dj
                if 0 <= ni < H and 0 <= nj < W:
                    adj[pos].append(ni * W + nj)

    from functools import lru_cache

    # 使用最小可选数优先的顺序
    # 预先计算每个位置的可能颜色
    possible = [
        set([1, 2, 3]) if flat_grid[pos] == -1 else set([flat_grid[pos]]) for pos in range(total)
    ]

    # 排序位置，优先填充约束更多的区域
    positions = list(range(total))
    positions.sort(key=lambda x: len(adj[x]))

    @lru_cache(maxsize=None)
    def dfs(index, assignment):
        if index == total:
            return 1

        pos = positions[index]
        if flat_grid[pos] != -1:
            return dfs(index + 1, assignment)

        # 获取相邻已赋值的颜色
        forbidden = set()
        for neighbor in adj[pos]:
            color = (assignment >> (neighbor * 2)) & 3
            if color != 0:
                forbidden.add(color)

        ans = 0
        for color in range(1, 4):
            if color not in forbidden:
                # 设置当前颜色
                new_assignment = assignment | (color << (pos * 2))
                ans = (ans + dfs(index + 1, new_assignment)) % MOD
        return ans

    # 初始赋值
    initial_assignment = 0
    for pos in range(total):
        if flat_grid[pos] != -1:
            color = flat_grid[pos]
            initial_assignment |= color << (pos * 2)

    result = dfs(0, initial_assignment)
    print(result)


if __name__ == "__main__":
    solve()
