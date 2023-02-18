# 对于每个点的度数都不超过2的无向图，
# 可以将路径的顶点序列和环的顶点序列分解出来。
# 返回：路径的组, 环的组

from typing import List, Tuple


def pathCycle(n: int, adjList: List[List[int]]) -> Tuple[List[List[int]], List[List[int]]]:
    if n == 0:
        return [], []
    deg = [len(nexts) for nexts in adjList]
    assert max(deg) <= 2, "degree of each vertex must not exceed 2"

    def calFrom(v: int) -> List[int]:
        path = [v]
        visited[v] = True
        while True:
            ok = False
            for next in adjList[path[-1]]:
                if visited[next]:
                    continue
                path.append(next)
                visited[next] = True
                ok = True
            if not ok:
                break
        return path

    visited = [False] * n
    paths, cycles = [], []
    for v in range(n):
        if deg[v] == 0:
            visited[v] = True
            paths.append([v])
        if visited[v] or deg[v] != 1:
            continue
        paths.append(calFrom(v))
    for v in range(n):
        if visited[v]:
            continue
        cycles.append(calFrom(v))
    return paths, cycles
