"""
栈实现dfs
适用于atcoder
"""

n = 5
edges = [[0, 1], [0, 2], [1, 3], [1, 4]]
adjList = [[] for _ in range(n)]
for u, v in edges:
    adjList[u].append(v)
    adjList[v].append(u)


def dfs(root: int) -> None:
    """栈实现dfs"""
    stack = [(root, -1, 0)]  # cur, pre, dep
    path = [0] * n  # 记录路径(每个深度对应的结点)
    while stack:
        cur, pre, dep = stack.pop()
        path[dep] = cur
        # print(dep, cur)  # ...处理当前结点的逻辑
        for next in adjList[cur]:
            if next == pre:
                continue
            stack.append((next, cur, dep + 1))


dfs(0)
