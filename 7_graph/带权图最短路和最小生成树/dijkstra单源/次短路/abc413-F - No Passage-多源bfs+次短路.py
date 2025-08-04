# F - No Passage(多源bfs+次短路)
# 在 H×W 的网格中，有一枚棋子从任意格子 (i,j) 出发。操作流程如下：
#
# 青木选一个方向 a∈{1,2,3,4}（上、下、左、右）。
# 高桥再选一个不同方向 b≠a，并朝 b 方向移动棋子；若越界则原地不动。
# 这一轮算一次“移动”。这一过程会反复进行，直至棋子抵达某个目标格或永不抵达。
# 高桥想以最少的“移动”次数到达任一目标，青木想阻止或（若无法阻止）使高桥所需的“移动”次数尽可能多。
# 对每个起点格 (i,j)，若高桥可被逼到目标，其最少移动次数为 f(i,j)；否则结果 0。求所有 f(i,j) 的总和。
#
# !维护每个格 u 的两个最小已知邻居距离

from collections import deque


INF = int(4e18)
DIR4 = [(0, -1), (0, 1), (-1, 0), (1, 0)]

if __name__ == "__main__":
    H, W, K = map(int, input().split())
    targets = [tuple(map(int, input().split())) for _ in range(K)]
    for i in range(K):
        targets[i] = (targets[i][0] - 1, targets[i][1] - 1)

    dist = [[INF] * W for _ in range(H)]
    counter = [[0] * W for _ in range(H)]
    queue = deque()
    for x, y in targets:
        dist[x][y] = 0
        counter[x][y] = 2
        queue.append((x, y))

    while queue:
        r, c = queue.popleft()
        for dx, dy in DIR4:
            nr, nc = r + dx, c + dy
            if 0 <= nr < H and 0 <= nc < W:
                counter[nr][nc] += 1
                if counter[nr][nc] == 2:
                    dist[nr][nc] = dist[r][c] + 1
                    queue.append((nr, nc))

    res = sum(sum(v for v in row if v < INF) for row in dist)
    print(res)
