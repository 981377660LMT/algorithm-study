# 每个人都有讨厌的一个人
# 从这些人中选最多的点，使得相邻的点不能被同时选，问选的点权最大和是多少

# 给定一个 n 个点 n 条边的没有自环的图，每个点有权值，
# 有边相连的两个点只能选其一，求可选方案的最大的点权之和。
# n≤106

# 对每一棵基环树断掉环上的一条边，然后对断开的两个点分别跑树形DP，就可以得到正确的答案。

# 没有上司的舞会 带环版本
from collections import defaultdict, deque
from typing import DefaultDict, List, Set, Tuple
from 基环树找到所有环 import cyclePartition

AdjMap = DefaultDict[int, Set[int]]
Degrees = List[int]


def dfs(cur: int, removed: int) -> List[int]:
    """从环上选出若干不相邻的点使权重和最大

    后序dfs返回[不选当前，选当前]
    """
    res = [0, values[cur]]

    for next in radjMap[cur]:  # 注意这里用外向基环树处理子树
        if next == removed:
            continue
        noNext, hasNext = dfs(next, removed)
        res[0] += max(noNext, hasNext)
        res[1] += noNext

    return res


def main(n: int, adjMap: AdjMap) -> int:
    cycleGroup = cyclePartition(n, adjMap, directed=True)[0]  # 找到所有环分组
    res = 0
    # 从所有环开始dp
    for group in cycleGroup:
        if len(group) < 2:
            continue
        # 取环上相邻的两个点 分别以这两个点为根，求 max(f[root1][不选],f[root2][不选]) 即可
        # 只要我们不选当前点，自然不会和另一个点冲突，就可以把对应那条边断开
        root1, root2 = (group[0], group[1])
        res += max(dfs(root1, root1)[0], dfs(root2, root2)[0])  # 断开这条边，即从自己出发不能走回自己

    return res


n = int(input())
adjMap = defaultdict(set)  # 内向基环树
radjMap = defaultdict(set)  # 外向基环树
values = [0] * n

for cur in range(n):
    value, hate = map(int, input().split())
    hate -= 1
    values[cur] = value
    adjMap[cur].add(hate)
    radjMap[hate].add(cur)

print(main(n, adjMap))

# 如果是一个 n 个点 n 条边的连通图的话，那么就是一棵基环树
# 但是讨厌的关系不保证连通，所以是很多个基环树
# 我们对于每个基环树单独求最优解，加起来即可
# 5
# 799600 3
# 551723 5
# 704668 5
# 849165 3
# 367471 4
