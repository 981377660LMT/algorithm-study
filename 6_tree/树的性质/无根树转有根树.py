# 无根树(无向边)转有根树(有向边)


from collections import deque
from typing import List, Tuple


def toRootedTree(tree: List[List[int]], root=0) -> List[List[int]]:
    n = len(tree)
    res = [[] for _ in range(n)]
    visited = [False] * n
    visited[root] = True
    queue = deque([root])
    while queue:
        cur = queue.popleft()
        for next in tree[cur]:
            if not visited[next]:
                visited[next] = True
                queue.append(next)
                res[cur].append(next)
    return res


def toRootedTreeWeighted(tree: List[List[Tuple[int, int]]], root=0) -> List[List[Tuple[int, int]]]:
    n = len(tree)
    res = [[] for _ in range(n)]
    visited = [False] * n
    visited[root] = True
    queue = deque([root])
    while queue:
        cur = queue.popleft()
        for next, weight in tree[cur]:
            if not visited[next]:
                visited[next] = True
                queue.append(next)
                res[cur].append((next, weight))
    return res
