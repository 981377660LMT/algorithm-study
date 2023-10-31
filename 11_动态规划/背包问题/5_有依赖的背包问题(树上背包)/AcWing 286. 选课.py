# 有依赖的背包问题
# 注意这里是很多棵树 我们可以加一个虚拟结点0连接所有的树
# 为自己确定一个选课方案，选m门课，使得你能得到的学分最多，并且必须满足先修条件。
# 课程编号1-n
# 时间复杂度 O(n3)

from collections import defaultdict
from functools import lru_cache

n, m = map(int, input().split())
nodes = [(0, 0)]  # 虚拟节点体积为0，价值为0
adjMap = defaultdict(list)
for i in range(n):
    pre, value = map(int, input().split())
    nodes.append((1, value))  # 体积为1，价值为value
    adjMap[pre].append(i + 1)  # 如果pre为0时表示虚拟节点


@lru_cache(None)
def dfs(root: int, select: int) -> int:
    """root子树(包含root)里选select个课的最大价值"""
    if nodes[root][0] > select:
        return 0

    rootValue = nodes[root][1]
    select -= nodes[root][0]

    dp = [0] * (select + 1)
    # 分组背包：物品，容量，决策
    for i in range(len(adjMap[root])):
        for j in range(select, -1, -1):
            nextCost = nodes[adjMap[root][i]][0]
            for k in range(nextCost, j + 1):
                dp[j] = max(dp[j], dfs(adjMap[root][i], k) + dp[j - k])
    return rootValue + dp[select]


print(dfs(0, m))
