import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 頂点
# M 辺の連結な無向グラフがあり、
# i 番目の辺は頂点
# U
# i
# ​
#   と頂点
# V
# i
# ​
#   を双方向に結びます。
# また、全ての頂点に整数が書いてあり、頂点
# v には整数
# A
# v
# ​
#   が書かれています。

# 頂点
# 1 から頂点
# N への単純なパス ( 同じ頂点を複数回通らないパス ) に対して、以下のように得点を定めます。

# パス上の頂点に書かれた整数を通った順に並べた数列 を
# S とする。
# S が広義単調増加になっていない場合、そのパスの得点は
# 0 である。
# そうでない場合、
# S に含まれる整数の種類数が得点となる。
# 頂点
# 1 から頂点
# N への全ての単純なパスのうち、最も得点が高いものを求めてその得点を出力してください。
from collections import deque
from functools import lru_cache
from typing import Callable, List, Optional, Tuple


def findSCC(n: int, graph: List[List[int]]) -> Tuple[List[List[int]], List[int]]:
    """
    # !Tarjan 算法求有向图的 scc

    Args:
        n (int): 图的顶点数
        graph (List[List[int]]):  邻接表

    Returns:
        Tuple[List[List[int]], List[int]]:
        每个 scc 组里包含的点，每个点所在 scc 的编号(0 ~ len(groups)-1)
    """

    def dfs(cur: int) -> int:
        nonlocal dfsId
        dfsId += 1
        dfsOrder[cur] = dfsId
        curLow = dfsId
        stack.append(cur)
        inStack[cur] = True
        for next in graph[cur]:
            if dfsOrder[next] == 0:
                nextLow = dfs(next)
                if nextLow < curLow:
                    curLow = nextLow
            elif inStack[next] and dfsOrder[next] < curLow:
                curLow = dfsOrder[next]
        if dfsOrder[cur] == curLow:
            group = []
            while True:
                top = stack.pop()
                inStack[top] = False
                group.append(top)
                if top == cur:
                    break
            groups.append(group)
        return curLow

    dfsOrder = [0] * n
    dfsId = 0
    stack = []
    inStack = [False] * n
    groups = []
    for i, order in enumerate(dfsOrder):
        if order == 0:
            dfs(i)

    # 由于每个强连通分量都是在它的所有后继强连通分量被求出之后求得的
    # 上面得到的 scc 是拓扑序的逆序
    groups.reverse()
    belong = [0] * n
    for i, group in enumerate(groups):
        for v in group:
            belong[v] = i

    return groups, belong  # !groups按照拓扑序输出


# TODO: deg 有问题
def toDAG(
    graph: List[List[int]],
    groups: List[List[int]],
    sccId: List[int],
    f: Optional[Callable[[int, int, int, int], None]] = None,
) -> Tuple[List[List[int]], List[int]]:
    """
    !scc 缩点成DAG

    Args:
        - graph (List[List[int]]):  邻接表
        - groups (List[List[int]]): 每个 scc 组里包含的点
        - sccId (List[int]): 每个点所在 scc 的编号(从0开始)
        - f (Optional[Callable[[int, int, int, int], None]]):
          回调函数，入参为 `(from, fromSccId, to, toSccId)`.

    Returns:
        - dag: 缩点成DAG后的邻接表.
        - indeg: 缩点后每个点的入度.
    """
    m = len(groups)
    dag = [[] for _ in range(m)]
    visitedEdge = set()  # !去除重边
    indeg = [0] * m
    for cur, nexts in enumerate(graph):
        curId = sccId[cur]
        for next in nexts:
            nextId = sccId[next]
            if curId != nextId:
                hash_ = curId * m + nextId
                if hash_ in visitedEdge:
                    continue
                visitedEdge.add(hash_)
                dag[curId].append(nextId)
                indeg[nextId] += 1
            if f is not None:
                f(cur, curId, next, nextId)
    return dag, indeg


if __name__ == "__main__":
    n, m = map(int, input().split())
    if n == 1:
        print(1)
        exit()
    values = list(map(int, input().split()))
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        w1, w2 = values[u], values[v]
        if w1 < w2:
            adjList[u].append(v)
        elif w1 > w2:
            adjList[v].append(u)
        else:
            adjList[u].append(v)
            adjList[v].append(u)

    groups, sccId = findSCC(n, adjList)
    dag, indeg = toDAG(adjList, groups, sccId)

    # dp = [0] * len(groups)
    # zeroId = sccId[0]
    # queue = deque([zeroId])
    # dp[zeroId] = 1

    # while queue:
    #     cur = queue.popleft()
    #     for next in dag[cur]:
    #         dp[next] = max(dp[next], dp[cur] + 1)
    #         indeg[next] -= 1
    #         if indeg[next] == 0:
    #             queue.append(next)

    # print(dp[sccId[n - 1]])
    start, target = sccId[0], sccId[n - 1]

    @lru_cache(None)
    def dfs(cur: int) -> int:
        if cur == target:
            return 0
        res = -INF
        for next in dag[cur]:
            res = max(res, dfs(next) + 1)
        return res

    res = dfs(start)
    if res < 0:
        print(0)
    else:
        print(res + 1)
