"""
给定一个 n 个点 m 条边的无向图, 每个点都有一个权值 A_{i} , 
我们可以进行以下操作 n 次:
选择一个尚未被删除的点 x , 然后删除点 x 以及所有与之相连的点, 
代价是所有与之相连的未被删除的点的权值之和;
我们定义 n 次操作的代价为所有操作中代价的最大值, 找到 n 次操作的代价的最小值

二分+类似拓扑排序的过程
出队后更新一遍相邻边的cost 注意有环 要用visited数组
"""

from collections import deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, m = map(int, input().split())
values = list(map(int, input().split()))
adjList = [[] for _ in range(n)]
cost = [0] * n
for _ in range(m):
    u, v = map(int, input().split())
    u, v = u - 1, v - 1
    adjList[u].append(v)
    adjList[v].append(u)
    cost[u] += values[v]
    cost[v] += values[u]


def check(mid: int) -> bool:
    """n次操作的代价是否可以不超过mid"""
    curCost = cost[:]
    queue = deque([i for i in range(n) if curCost[i] <= mid])
    visited = [curCost[i] <= mid for i in range(n)]
    while queue:
        cur = queue.popleft()
        for next in adjList[cur]:
            if visited[next]:
                continue
            curCost[next] -= values[cur]
            if curCost[next] <= mid:
                queue.append(next)
                visited[next] = True
    return all(visited)


left, right = 0, int(1e9) * n
while left <= right:
    mid = (left + right) // 2
    if check(mid):
        right = mid - 1
    else:
        left = mid + 1
print(left)
