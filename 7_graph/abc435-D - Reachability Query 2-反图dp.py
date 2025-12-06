# 染色+可到达性查询(反图dp)
# https://atcoder.jp/contests/abc435/tasks/abc435_d
# !给定一个 $N$ 个顶点 $M$ 条边的有向图。 顶点编号为 $1$ 到 $N$，第 $i$ 条边是从顶点 $X_i$ 到顶点 $Y_i$ 的有向边。
# 最初所有顶点都是白色的。
# 给定 $Q$ 个查询，请按顺序处理。查询有以下两种类型：
#
# 1 v：将顶点 $v$ 染成黑色。
# 2 v：判断从顶点 $v$ 出发，沿着有向边是否能到达任意一个黑色的顶点。


from collections import deque


N, M = map(int, input().split())
revAdjList = [[] for _ in range(N)]
for _ in range(M):
    u, v = map(int, input().split())
    revAdjList[v - 1].append(u - 1)

visited = [False] * N  # 在反图中，从某个黑点出发能否到达 i
queue = deque()

Q = int(input())
for _ in range(Q):
    query = list(map(int, input().split()))
    kind, v = query
    v -= 1
    if kind == 1:
        # 染黑 v
        # 如果 v 已经在 visited 中，说明它已经被之前的黑点覆盖了，
        # 或者它自己之前就是黑点，无需重复搜索。
        if not visited[v]:
            visited[v] = True
            queue.append(v)
            while queue:
                cur = queue.popleft()
                for next in revAdjList[cur]:
                    if not visited[next]:
                        visited[next] = True
                        queue.append(next)

    else:
        print("Yes" if visited[v] else "No")
