# 它表示一个迷宫，其中的1表示墙壁，0表示可以走的路，只能横着走或竖着走，
# 不能斜着走，要求编程序找出从左上角到右下角的最短路线。

# 用一个pre数组存路径即可 第一次搜到的一定最短
from collections import deque


n = int(input())
mat = [list(map(int, input().split())) for _ in range(n)]

pre = [[None] * n for _ in range(n)]

queue = deque([(0, 0)])
while queue:
    r, c = queue.popleft()
    for dr, dc in zip((-1, 0, 1, 0), (0, -1, 0, 1)):
        nr, nc = r + dr, c + dc
        if 0 <= nr < n and 0 <= nc < n and mat[nr][nc] == 0 and pre[nr][nc] is None:
            pre[nr][nc] = (r, c)
            queue.append((nr, nc))

res = []
x, y = n - 1, n - 1
while (x, y) != (0, 0):
    res.append((x, y))
    x, y = pre[x][y]
res.append((0, 0))
for x, y in res[::-1]:
    print(x, y)
