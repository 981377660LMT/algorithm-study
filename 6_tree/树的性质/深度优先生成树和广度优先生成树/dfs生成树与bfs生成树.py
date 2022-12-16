# dfs生成树与bfs生成树
# https://www.cnblogs.com/miraclepbc/p/16280781.html

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
