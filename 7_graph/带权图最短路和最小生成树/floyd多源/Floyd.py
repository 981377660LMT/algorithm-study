from collections import defaultdict
from itertools import product


INF = int(1e20)


if __name__ == "__main__":
    n = 3
    edges = [[0, 1, 2], [1, 0, 3], [2, 1, 1]]  # 无向带权边
    dist = defaultdict(lambda: defaultdict(lambda: INF))
    for u, v, w in edges:
        dist[u][v] = min(dist[u][v], w)
        dist[v][u] = min(dist[v][u], w)
    for k, i, j in product(range(n), repeat=3):
        dist[i][j] = min(dist[i][j], dist[i][k] + dist[k][j])
    print(dist)
