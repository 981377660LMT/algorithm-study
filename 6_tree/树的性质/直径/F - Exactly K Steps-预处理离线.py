"""
给定一棵n(n<=2e5)个节点的树。
有q(q≤2e5)次询问,每次询问给出两个数字(u, k),
请找到距离点u的为k的点(输出任意一个即可)。如果没有输出-1。

q 个查询 询问距离树结点u距离为k的结点是否存在,并输出一个这样的结点

离线查询预处理
从直径端点出发dfs 记录好路径
每个结点处看是否命中查询
"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


from typing import List, Tuple
from collections import deque


def calDiameter1(adjList: List[List[int]], start: int) -> Tuple[int, Tuple[int, int]]:
    """bfs计算树的直径长度和直径两端点"""
    queue = deque([start])
    visited = set([start])
    last1 = 0  # 第一次BFS最后一个点
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            last1 = queue.popleft()
            for next in adjList[last1]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)

    queue = deque([last1])  # 第一次最后一个点作为第二次BFS的起点
    visited = set([last1])
    last2 = 0  # 第二次BFS最后一个点
    res = -1
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            last2 = queue.popleft()
            for next in adjList[last2]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)
        res += 1

    return res, tuple(sorted([last1, last2]))


if __name__ == "__main__":

    n = int(input())
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)

    _, (left, right) = calDiameter1(adjList, start=0)  # 求出直径两端点

    q = int(input())
    res = [-1] * q
    queries = [[] for _ in range(n)]  # !记录每个结点处的查询
    for i in range(q):
        u, k = map(int, input().split())  # 询问距离树结点u距离为k的结点是否存在,并输出一个这样的结点
        u = u - 1
        queries[u].append((i, k))

    # !TLE
    # def dfs(cur: int, pre: int, dep: int, path: List[int]) -> None:
    #     """从直径端点出发处理查询"""
    #     path[dep] = cur
    #     for i, k in queries[cur]:
    #         if dep - k >= 0:
    #             res[i] = path[dep - k]
    #     for next in adjList[cur]:
    #         if next == pre:
    #             continue
    #         dfs(next, cur, dep + 1, path)

    # dfs(left, -1, 0, [-1] * n)
    # dfs(right, -1, 0, [-1] * n)
    # print(*[num + 1 if num != -1 else -1 for num in res], sep="\n")

    def dfs(root: int) -> None:
        stack = [(root, -1, 0)]
        path = [-1] * n
        while stack:
            cur, pre, dep = stack.pop()
            path[dep] = cur
            for i, k in queries[cur]:
                if dep - k >= 0:
                    res[i] = path[dep - k]

            for next in adjList[cur]:
                if next == pre:
                    continue
                stack.append((next, cur, dep + 1))

    dfs(left)
    dfs(right)
    print(*[num + 1 if num != -1 else -1 for num in res], sep="\n")
