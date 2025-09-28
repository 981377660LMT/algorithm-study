# abc425-D - Ulam-Warburton Automaton-元胞自动机(生命游戏、乌拉姆-沃伯顿自动机)
# https://atcoder.jp/contests/abc425/tasks/abc425_d
# 有一个 H 行 W 列的网格。从上往下第 i 行（1 ≤ i ≤ H）、从左往右第 j 列（1 ≤ j ≤ W）的格子被称为格子 (i, j)。
#
# 初始时，如果 S_ij 是 #，格子 (i, j) 被涂成黑色；如果是 .，则被涂成白色。
#
# 接下来，执行 10^100 次以下操作：
#
# 令 T 为所有满足以下条件的白色格子的集合：该格子的四个邻边格子中，恰好只有一个是黑色的。
# 将集合 T 中的每一个格子都涂成黑色。
# （注：两个格子 (i1, j1) 和 (i2, j2) 邻边相邻，当且仅当 |i1 - i2| + |j1 - j2| = 1。）
#
# 请计算所有操作结束后，被涂成黑色的格子的总数。
#
# !每个格子只会被涂黑一次，故直接模拟每轮操作，直到没有格子可以被涂黑为止即可。


DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]

H, W = map(int, input().split())
grid = [input() for _ in range(H)]


visited = [[False] * W for _ in range(H)]
queue = []
for i in range(H):
    for j in range(W):
        if grid[i][j] == "#":
            visited[i][j] = True
            queue.append((i, j))

while queue:
    nextQueue = []
    cands = set()
    for r, c in queue:
        for dr, dc in DIR4:
            nr, nc = r + dr, c + dc
            if 0 <= nr < H and 0 <= nc < W and not visited[nr][nc]:
                cands.add((nr, nc))

    for r, c in cands:
        count = 0
        for dr, dc in DIR4:
            nr, nc = r + dr, c + dc
            if 0 <= nr < H and 0 <= nc < W and visited[nr][nc]:
                count += 1
        if count == 1:
            nextQueue.append((r, c))

    queue = nextQueue
    for r, c in queue:
        visited[r][c] = True

print(sum(sum(row) for row in visited))
