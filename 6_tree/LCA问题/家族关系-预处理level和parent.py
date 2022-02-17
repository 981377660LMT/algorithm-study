# 亲疏关系可以使用数字表示，如父子关系是 1，爷孙关系和兄弟关系是 2，以此类推。
# 辈份关系是指第一个人高于第二个人的辈份数，可以是正数，负数或 0。


# 注意这道题坑的地方：根节点不一定是0
# 1.建图，建图的过程中记录每个点的入度，用于寻找根节点
# 2.从根节点出发dfs或者bfs预处理 parent 和 level
# 3.根据 parent 和 level 找lca

from collections import defaultdict
from typing import DefaultDict


def dfs(cur: int, parent: int, depth: int) -> None:
    parentMap[cur] = parent
    levelMap[cur] = depth
    for next in adjMap[cur]:
        if next == parent:
            continue
        dfs(next, cur, depth + 1)


def LCA(root1: int, root2: int, level: DefaultDict[int, int], parent: DefaultDict[int, int]) -> int:
    if level[root1] < level[root2]:
        root1, root2 = root2, root1
    diff = level[root1] - level[root2]
    for _ in range(diff):
        root1 = parent[root1]
    while root1 != root2:
        root1 = parent[root1]
        root2 = parent[root2]
    return root1


n, root1, root2 = map(int, input().split())
adjMap = defaultdict(list)
indegree = [0] * n  # 记录每个点的入度，用于寻找根节点
for i in range(n - 1):
    u, v = map(int, input().split())
    adjMap[u].append(v)
    adjMap[v].append(u)
    indegree[v] += 1

levelMap, parentMap = defaultdict(lambda: -1), defaultdict(lambda: -1)
root = next(i for i in range(n) if indegree[i] == 0)
dfs(root, -1, 0)

lca = LCA(root1, root2, levelMap, parentMap)


# 求树中两点距离以及高度差
print(
    levelMap[root1] - levelMap[lca] + levelMap[root2] - levelMap[lca],
    levelMap[root2] - levelMap[root1],
)
