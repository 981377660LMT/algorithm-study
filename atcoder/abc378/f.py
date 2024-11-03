from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    N = int(sys.stdin.readline())
    adj = [[] for _ in range(N)]
    degrees = [0] * N
    edges = set()
    for _ in range(N - 1):
        u, v = map(int, sys.stdin.readline().split())
        u -= 1
        v -= 1
        adj[u].append(v)
        adj[v].append(u)
        degrees[u] += 1
        degrees[v] += 1
        edges.add((min(u, v), max(u, v)))

    degree3_nodes = [i for i in range(N) if degrees[i] == 3]

    visited = [False] * N

    components = []

    def dfs(u, component_nodes, degree2_nodes):
        visited[u] = True
        component_nodes.add(u)
        for v in adj[u]:
            if degrees[v] == 3 and not visited[v]:
                dfs(v, component_nodes, degree2_nodes)
            elif degrees[v] == 2:
                degree2_nodes.add(v)

    for node in degree3_nodes:
        if not visited[node]:
            component_nodes = set()
            degree2_nodes = set()
            dfs(node, component_nodes, degree2_nodes)
            components.append((component_nodes, degree2_nodes))

    total_pairs = 0
    for component_nodes, degree2_nodes in components:
        m = len(degree2_nodes)
        total_pairs += m * (m - 1) // 2

    print(total_pairs)
