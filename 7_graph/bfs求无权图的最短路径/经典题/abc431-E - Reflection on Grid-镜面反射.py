# abc431-E - Reflection on Grid-镜面反射/镜子反射
# https://atcoder.jp/contests/abc431/tasks/abc431_e
#
# 有一个H行W列的网格。
# 每个格子至多放置一面镜子。高桥站在格子(1,1)的左侧，青木站在格子(H,W)的右侧。
# 高桥有一支手电筒，并从格子(1,1)左侧向右照射光线。
# 假设手电筒发出的光不会发散，而是一条沿直线传播的光线。
# 高桥的目标是利用网格中的镜子，将手电筒的光线传递给青木。
# 镜子有三种放置类型。当光线照射到镜子时，其传播方向会根据镜子的类型发生改变。
# A：无镜子
# B：\ 镜子
# C：/ 镜子
# 高桥可以执行以下操作任意多次，以将光线传递给青木：
# 选择一个格子，将其镜子类型更改为另一种类型。
# 求将光线成功传递给青木所需的最少操作次数。
#
# 01bfs
#
# 更换一个格子的镜子相当于改变光线在此格发射出去的方向，
# 因此可以：
# - 将光线正常传播视为不花费操作；
# - 将更换镜子视为花费一次操作改变光线在此格发射出去的方向。

from collections import deque


INF = int(1e18)
DIR4 = [(-1, 0), (0, 1), (1, 0), (0, -1)]  # 上右下左


def mirror(c: str, d: int) -> int:
    if c == "A":
        return d
    if c == "B":
        return 3 - d
    return d ^ 1


def solve(H: int, W: int, grid: list[str]) -> int:
    d = mirror(grid[0][0], 1)
    dist = [[[INF] * 4 for _ in range(W)] for _ in range(H)]
    dist[0][0][d] = 0
    queue = deque([(0, 0, d)])
    while queue:
        r, c, d = queue.popleft()
        nr, nc = r + DIR4[d][0], c + DIR4[d][1]
        if 0 <= nr < H and 0 <= nc < W:
            nd = mirror(grid[nr][nc], d)
            if dist[nr][nc][nd] > dist[r][c][d]:
                dist[nr][nc][nd] = dist[r][c][d]
                queue.appendleft((nr, nc, nd))
        for nd in range(4):
            if nd == d:
                continue
            if dist[r][c][nd] > dist[r][c][d] + 1:
                dist[r][c][nd] = dist[r][c][d] + 1
                queue.append((r, c, nd))
    return dist[H - 1][W - 1][1]


for _ in range(int(input())):
    H, W = map(int, input().split())
    grid = [input() for _ in range(H)]
    res = solve(H, W, grid)
    print(res)
