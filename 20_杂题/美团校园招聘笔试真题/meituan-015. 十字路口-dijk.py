from heapq import heappop, heappush
from collections import defaultdict

INF = 0x3F3F3F3F
# dijk
n, m, xs, ys, xt, yt = map(int, input().split())
a, b = [], []
for i in range(n):
    a.append(list(map(int, input().split())))
for i in range(n):
    b.append(list(map(int, input().split())))


queue = [(0, xs - 1, ys - 1)]
dist = defaultdict(lambda: INF)
dist[(xs - 1, ys - 1)] = 0

while queue:
    t, x, y = heappop(queue)
    if t % (a[x][y] + b[x][y]) < a[x][y]:
        if x + 1 < n and dist[(x + 1, y)] > t + 1:
            dist[(x + 1, y)] = t + 1
            heappush(queue, [t + 1, x + 1, y])

        if x - 1 >= 0 and dist[(x - 1, y)] > t + 1:
            dist[(x - 1, y)] = t + 1
            heappush(queue, [t + 1, x - 1, y])

    elif b[x][y] + a[x][y] > t % (a[x][y] + b[x][y]) >= a[x][y]:
        if y + 1 < m and dist[(x, y + 1)] > t + 1:
            dist[(x, y + 1)] = t + 1
            heappush(queue, [t + 1, x, y + 1])

        if y - 1 >= 0 and dist[(x, y - 1)] > t + 1:
            dist[(x, y - 1)] = t + 1
            heappush(queue, [t + 1, x, y - 1])

    cur = t % (a[x][y] + b[x][y])
    if cur < a[x][y]:
        cost = a[x][y] - cur + 1
        if y + 1 < m and dist[(x, y + 1)] > cost + t:
            dist[(x, y + 1)] = cost + t
            heappush(queue, [cost + t, x, y + 1])
        if y - 1 >= 0 and dist[(x, y - 1)] > cost + t:
            dist[(x, y - 1)] = cost + t
            heappush(queue, [cost + t, x, y - 1])

    elif cur < a[x][y] + b[x][y]:
        cost = a[x][y] + b[x][y] - cur + 1
        if x + 1 < n and dist[(x + 1, y)] > cost + t:
            dist[(x + 1, y)] = cost + t
            heappush(queue, [cost + t, x + 1, y])

        if x - 1 >= 0 and dist[(x - 1, y)] > cost + t:
            dist[(x - 1, y)] = cost + t
            heappush(queue, [cost + t, x - 1, y])

print(dist[(xt - 1, yt - 1)])

