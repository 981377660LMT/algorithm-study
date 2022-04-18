from collections import defaultdict

# 2 <= n <= 1000,1 <= m <= 1000,1 <= q <= 1000
n, m = map(int, input().split())

adjMatrix = defaultdict(lambda: defaultdict(lambda: int(1e20)))

for _ in range(m):
    u, v, w = input().split()
    w = int(w)
    adjMatrix[u][v] = min(adjMatrix[u][v], w)

keys = list(adjMatrix.keys())
for k in keys:
    for i in keys:
        for j in keys:
            adjMatrix[i][j] = min(adjMatrix[i][j], adjMatrix[i][k] + adjMatrix[k][j])


q = int(input())
for _ in range(q):
    u, v = input().split()
    res = adjMatrix[u][v]
    print(res if res < int(1e19) else 'INF')
