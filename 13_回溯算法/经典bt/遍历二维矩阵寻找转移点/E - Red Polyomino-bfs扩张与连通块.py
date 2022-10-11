"""
bfs+状态压缩
"""

# 红色骨牌(七巧板)

# 给你 一个 n*n的矩阵,
# 其中各单元属性要么是白要么是黑,
# 现在需要你计算出有多少种方案可以将`白块绘制成红块`,
# !使红块区域连通且大小为k。
# !n,k<=8
# !"#"表示黑色,"."表示白色

# 1. comb(64,8)>4e9 愚直な全探索は間に合いません。
# !2. しかし、`赤マス同士が連結`という条件から、
# 条件を満たす選び方はそれより相当に少ないことが予想できます
# 题目的最坏样例告诉我们,这样的情况是很少的

# !直接bfs+visited 存储看过的状态来搜索
# !(bfs扩张k层,就可以得到大小为k的连通块)
# 64也在暗示状态压缩

from collections import deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


DIR4 = [(1, 0), (-1, 0), (0, 1), (0, -1)]


if __name__ == "__main__":
    n = int(input())
    k = int(input())
    grid = [list(input()) for _ in range(n)]

    visited = set()  # 用于存储访问过的红色连通块状态
    queue = deque()
    for r in range(n):
        for c in range(n):
            if grid[r][c] == ".":
                pos = r * n + c
                queue.append(1 << pos)
                visited.add(1 << pos)

    for _ in range(k - 1):  # bfs扩张k-1层
        len_ = len(queue)
        for _ in range(len_):
            curState = queue.popleft()
            #  !遍历整个矩阵,寻找下个红色块,判断和之前的红色块是否连通
            for r in range(n):
                for c in range(n):
                    if grid[r][c] == "#" or curState & (1 << (r * n + c)):
                        continue

                    ok = False
                    for dr, dc in DIR4:
                        nr, nc = r + dr, c + dc
                        if 0 <= nr < n and 0 <= nc < n and curState & (1 << (nr * n + nc)):
                            ok = True
                            break
                    if ok:
                        nextState = curState | (1 << (r * n + c))
                        if nextState not in visited:
                            visited.add(nextState)
                            queue.append(nextState)

    print(len(queue))
