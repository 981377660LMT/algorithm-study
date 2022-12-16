"""
无向图G中求两个生成树T1和T2
当1为根时
要求G中不在T1里任意两条边都有祖孙关系
要求G中不在T2里任意两条边都没有祖孙关系
输出这两棵生成树的边
# n,m<=2e5

利用dfs与bfs的性质

生成树:
连接图中所有的点n,并且只有n-1条边的子图就是它的生成树
只要能连通所有顶点而又不产生回路的任何子图都是它的生成树
"""


from collections import deque
from typing import List, Tuple


def dfsSpanningTree(n: int, adjList: List[List[int]]) -> List[Tuple[int, int]]:
    """dfs生成树

    - 不在dfs生成树中的任意一条边的两个顶点在生成树中都`有`祖孙关系。
    """

    def dfs(cur: int) -> None:
        for next in adjList[cur]:
            if not visited[next]:
                visited[next] = True
                edges.append((cur, next))
                dfs(next)

    visited = [False] * n
    visited[0] = True
    edges = []
    dfs(0)
    return edges


def bfsSpanningTree(n: int, adjList: List[List[int]]) -> List[Tuple[int, int]]:
    """bfs生成树

    - 不在bfs生成树中的任意一条边的两个顶点在生成树中都`没有`祖孙关系。
    - bfs生成树在图的所有生成树中的`高度是最低的`。
    """
    visited = [False] * n
    visited[0] = True
    queue = deque([0])
    edges = []
    while queue:
        cur = queue.popleft()
        for next in adjList[cur]:
            if not visited[next]:
                visited[next] = True
                queue.append(next)
                edges.append((cur, next))
    return edges


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)

    dfsEdges = dfsSpanningTree(n, adjList)
    for u, v in dfsEdges:
        print(u + 1, v + 1)
    bfsEdges = bfsSpanningTree(n, adjList)
    for u, v in bfsEdges:
        print(u + 1, v + 1)
