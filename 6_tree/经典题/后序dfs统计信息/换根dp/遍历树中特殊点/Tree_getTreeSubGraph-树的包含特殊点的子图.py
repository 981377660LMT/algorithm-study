# Tree_getTreeSubGraph-树中包含特殊点的导出子图
from typing import List, Tuple

INF = int(1e18)


def getTreeSubGraph(
    n: int, rawTree: List[List[int]], specials: List[int]
) -> Tuple[List[List[int]], List[Tuple[int, int]], List[bool]]:
    """给定`原树`(无向图邻接表), 返回`子图(无向图邻接表), 子图中的边, 原图中每个点是否在子图中`."""
    if len(specials) == 0:
        return [[] for _ in range(n)], [], [False] * n
    visited = [False] * n
    for v in specials:
        visited[v] = True

    def dfs(cur: int, pre: int) -> bool:
        for next in rawTree[cur]:
            if next != pre and dfs(next, cur):
                visited[cur] = True  # 标记必经点
        return visited[cur]

    root = specials[0]
    dfs(root, -1)

    edges, newTree = [], [[] for _ in range(n)]
    for cur in range(n):
        for next in rawTree[cur]:
            if cur < next and visited[cur] and visited[next]:
                edges.append((cur, next))
                newTree[cur].append(next)
                newTree[next].append(cur)
    return newTree, edges, visited
