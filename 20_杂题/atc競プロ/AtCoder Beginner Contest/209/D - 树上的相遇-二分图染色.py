# 树上的相遇
# 给定一棵无向无权的树和q个询问，每个询问给定两个点u,v
# 问u,v是在点相遇还是在边相遇(即u,v间路径长度是奇数还是偶数)

# 解法1:LCA求路径长度
# 解法2:树的二分图染色 O(n)

# Collision

from collections import deque
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def getTreeColor(n: int, adjList: List[List[int]]) -> List[int]:
    """树的二分图染色"""
    colors = [-1] * n
    colors[0] = 0
    queue = deque([0])
    while queue:
        cur = queue.popleft()
        for next in adjList[cur]:
            if colors[next] == -1:
                colors[next] = 1 ^ colors[cur]
                queue.append(next)
    return colors


if __name__ == "__main__":
    n, q = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)

    treeColor = getTreeColor(n, adjList)
    for _ in range(q):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        if treeColor[u] == treeColor[v]:
            print("Town")
        else:
            print("Road")
