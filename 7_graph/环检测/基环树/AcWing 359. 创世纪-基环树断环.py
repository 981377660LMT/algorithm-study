# 注意到题意中给出的每个点都能够限制某个点，如果从图论角度考虑，那么可以想到基环树。
# https://www.acwing.com/activity/content/problem/content/3188/


# 上帝手中有 N 种世界元素，每种元素可以限制另外1种元素，把第 i 种世界元素能够限制的那种世界元素记为 A[i]。

# 现在，上帝要把它们中的一部分投放到一个新的空间中去建造世界。

# 为了世界的和平与安宁，上帝希望所有被投放的世界元素都有至少一个没有被投放的世界元素限制它。

# 上帝希望知道，在此前提下，他最多可以投放多少种世界元素？


from collections import defaultdict, deque
from typing import DefaultDict, List, Set, Tuple
from 基环树找到所有环 import cyclePartition

AdjMap = DefaultDict[int, Set[int]]
Degrees = List[int]


def dfs(cur: int, removed: int) -> List[int]:
    """从环上选出若干不相邻的点使权重和最大

    后序dfs返回[不选当前，选当前]
    """

    res = [0, -int(1e20)]
    toDel = int(1e20)

    for next in radjMap[cur]:  # 注意这里用外向基环树处理子树
        if next == removed:
            continue
        noNext, hasNext = dfs(next, removed)
        res[0] += max(noNext, hasNext)  # 子节点怎么选都无所谓，直接都选最大的
        toDel = min(toDel, max(noNext, hasNext)) - noNext
    res[1] = res[0] + 1 - toDel  # 至少一个子节点需要不选取，所以我们贪心地决策：求取在一个子节点强制不选取的前提下的最小损失 del 即可。
    if removed in adjMap[cur]:
        res[1] = res[0] + 1
    return res


def main(n: int, adjMap: AdjMap) -> int:
    cycleGroup = cyclePartition(n, adjMap, directed=True)[0]  # 找到所有环分组
    res = 0
    # 从所有环开始dp
    for group in cycleGroup:
        # 取环上相邻的两个点 分别以这两个点为根，求 max(f[root1][不选],f[root2][不选]) 即可
        # 只要我们不选当前点，自然不会和另一个点冲突，就可以把对应那条边断开
        root1, root2 = (group[0], group[1])
        res += max(dfs(root1, root1)[0], dfs(root2, root2)[0])  # 断开这条边，即从自己出发不能走回自己

    return res


n = int(input())
nums = list(map(int, input().split()))
adjMap = defaultdict(set)  # 内向基环树
radjMap = defaultdict(set)  # 外向基环树

for u, v in enumerate(nums):
    v -= 1
    adjMap[u].add(v)

    radjMap[v].add(u)


print(main(n, adjMap))

# 树形dp不太对
