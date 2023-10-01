"""求基环树的最长链"""
# 基环树的最长链有两种情况：
# !第一种是在某棵树里，不经过环，即直径
# !第二种是经过环，分解成三段，即depth[x]+depth[y]+distOncycle(x,y)
# 然后是环上求最值，破环成链，这个表达式最大值可以滑动窗口求最值由单调队列求出，参考AcWing 289. 环路运输

# 我们可以先用一次DFS找出环，在每个环上节点出发处理树上的最长路径，
# 并计算出从根节点最远到达节点距离根节点的距离。
# 显然，如果经过环的话，肯定要经过这个最大的距离。

# https://www.acwing.com/solution/content/63035/
from collections import defaultdict, deque
from heapq import nlargest
from typing import DefaultDict, List, Set, Tuple
from 基环树找到所有环 import cyclePartition

AdjMap = DefaultDict[int, DefaultDict[int, int]]
Degrees = List[int]


def calMax1(root: int, rAdjMap: AdjMap, cycle: Set[int]) -> int:
    """求子树直径的长度(有向图)"""

    def dfs(cur: int) -> int:
        nonlocal res
        nexts = [0, 0]
        for next, weight in rAdjMap[cur].items():
            if next in cycle:
                continue
            nexts.append(dfs(next) + weight)
        max1, max2 = nlargest(2, nexts)
        res = max(res, max1 + max2)
        return max1

    res = 0
    dfs(root)
    return res


def calMax2(scores: List[int], dists: List[int]) -> int:
    """环上求depth[x]+depth[y]+preSum[x]-preSum[y]的最大值

    环上每个点i求出i前面n−1个点中使得depth[j]-j最大的点(其实求出值即可)
    """
    n = len(scores)
    res = 0
    scores *= 2  # 破环成链
    preSum = [0]
    for i in range(n):
        preSum.append(preSum[-1] + dists[i])
    for i in range(n):
        preSum.append(preSum[-1] + dists[i])

    maxQueue = deque()
    for i in range(2 * n):
        while maxQueue and maxQueue[0][1] <= i - n:
            maxQueue.popleft()
        if maxQueue:
            res = max(res, scores[i] + preSum[i] + maxQueue[0][0])
        while maxQueue and maxQueue[-1][0] < (scores[i] - preSum[i]):
            maxQueue.pop()
        maxQueue.append([scores[i] - preSum[i], i])
    return res


def calLongestPath(adjMap: AdjMap) -> int:
    """求基环树最长链"""
    cycleGroup, *_, depth = cyclePartition(n, adjMap, directed=True)  # 所有在环上的点，所有点的最大深度
    onCycle = set([v for g in cycleGroup for v in g])

    res = 0
    for group in cycleGroup:  # 遍历所有的环 在每个环(联通分量)里找最长的链
        cand1 = 0  # 最长链在树里，不经过环
        cand2 = 0  # 最长链经过环上
        for root in group:
            cand1 = max(cand1, calMax1(root, rAdjMap, onCycle))  # !反图外向基环树处理子树

        scores, dists = (
            [depth[i] for i in group],
            [adjMap[u][v] for u, v in zip(group, group[1:] + [group[0]])],
        )

        cand2 = calMax2(scores, dists)

        res += max(cand1, cand2)

    return res


n = int(input())
adjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))  # !内向基环树用来求环分组
rAdjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))  # !外向基环树用来求子树链
# 岛屿 i 上建了一座通向岛屿 a 的桥，桥的长度为 L。
for u in range(n):
    v, w = map(int, input().split())
    v -= 1
    adjMap[u][v] = w
    rAdjMap[v][u] = w


print(calLongestPath(adjMap))
