from typing import List, Tuple


def getTreeDiameter(n: int, tree: List[List[int]], start=0) -> Tuple[int, List[int]]:
    """求无权树的(直径长度,直径路径)."""

    def dfs(start: int) -> Tuple[int, List[int]]:
        dist = [-1] * n
        dist[start] = 0
        stack = [start]
        while stack:
            cur = stack.pop()
            for next in tree[cur]:
                if dist[next] != -1:
                    continue
                dist[next] = dist[cur] + 1
                stack.append(next)
        endPoint = dist.index(max(dist))
        return endPoint, dist

    u, _ = dfs(start)
    v, dist = dfs(u)
    diameter = dist[v]
    path = [v]
    while u != v:
        for next in tree[v]:
            if dist[next] + 1 == dist[v]:
                path.append(next)
                v = next
                break

    return diameter, path
