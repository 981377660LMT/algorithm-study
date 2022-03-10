# 给定一棵树，树中包含 n 个结点（编号1~n）和 n−1 条无向边，每条边都有一个权值。
# 现在请你找到树中的一条最长路径。
# 1≤n≤10000 ,
from collections import defaultdict


n = int(input())
adjMap = defaultdict(set)

for _ in range(n - 1):
    u, v, w = map(int, input().split())
    adjMap[u].add((v, w))
    adjMap[v].add((u, w))

res = 0


def dfs(cur: int, parent: int) -> int:
    global res
    """后序dfs，返回每个root处的最大路径长度"""
    # 每个点处最长和次长
    max1, max2 = 0, 0
    for next, weight in adjMap[cur]:
        if next == parent:
            continue
        maxCand = dfs(next, cur) + weight
        if maxCand > max1:
            max2, max1 = max1, maxCand
        elif maxCand > max2:
            max2 = maxCand
    res = max(res, max1 + max2)
    return max1


dfs(1, -1)
print(res)
