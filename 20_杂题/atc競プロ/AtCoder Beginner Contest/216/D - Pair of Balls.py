"""羊了个羊

有n个数,每个数都会出现两次。m个柱子,并给出柱子上的数。
当且仅当两个相同的数都处于柱子的最上面时,才能够将两个数删除。
问最后能否将柱子的数都删除。
n,m<=2e5

处于柱子的最上面:不存在依赖
!拓扑排序(注意到消除下面之前必须要消除上面)
"""


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

from typing import List
from collections import deque


def topoSort(n: int, adjList: List[List[int]], deg: List[int], directed: bool) -> List[int]:
    """求图的拓扑排序

    Args:
        n (int): 顶点0~n-1
        adjList (List[List[int]]): 邻接表
        deg (List[int]): 有向图的入度/无向图的度
        directed (bool): 是否为有向图

    Returns:
        List[int]: 拓扑排序结果, 若不存在则返回空列表
    """
    startDeg = 0 if directed else 1
    queue = deque([v for v in range(n) if deg[v] == startDeg])
    res = []
    while queue:
        cur = queue.popleft()
        res.append(cur)
        for next in adjList[cur]:
            deg[next] -= 1
            if deg[next] == startDeg:
                queue.append(next)

    return [] if len(res) < n else res


if __name__ == "__main__":
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    deg = [0] * n
    for _ in range(m):
        size = int(input())
        nums = [int(x) - 1 for x in map(int, input().split())]
        for pre, cur in zip(nums, nums[1:]):
            adjList[pre].append(cur)
            deg[cur] += 1
    res = topoSort(n, adjList, deg, True)
    print("Yes" if len(res) == n else "No")
