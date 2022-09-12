from typing import List


def stackBfs(n: int, adjList: List[List[int]], root=0) -> List[int]:
    """返回bfs遍历顺序"""
    stack = [root]
    order = [root]
    parent = [-1] * n
    while stack:
        cur = stack.pop()
        for next in adjList[cur]:
            if next == parent[cur]:
                continue
            parent[next] = cur
            order.append(next)
            stack.append(next)
    return order


if __name__ == "__main__":
    n = 5
    edges = [[1, 2], [1, 3], [2, 4], [2, 5]]
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)
    print(stackBfs(n, adjList, root=0))
