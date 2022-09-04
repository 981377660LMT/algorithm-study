# 'D' 表示两个数字间的递减关系，'I' 表示两个数字间的递增关系
# 现在你的任务是找到具有最小字典序的 [1, 2, ... n] 的排列，使其能代表输入的 秘密签名。

# !按 'I' 分段，每一段尽量字典序最小即可。
# n<=1e5


from heapq import heappop, heappush
from itertools import pairwise, permutations
from typing import List


def findPermutation(pattern: str, *, isMin: bool) -> List[int]:
    """寻找DI序列的字典序最小/最大排列  O(n)"""
    if isMin:
        res = []
        min_ = 1
        for down in pattern.split("I"):
            count = len(down)
            res.extend(range(min_ + count, min_ - 1, -1))  # count+1 个数
            min_ += count + 1
        return res
    else:
        n = len(pattern) + 1
        res = []
        max_ = n
        for up in pattern.split("D"):
            count = len(up)
            res.extend(range(max_ - count, max_ + 1))
            max_ -= count + 1
        return res


print(findPermutation("DI", isMin=True))
print(findPermutation("DI", isMin=False))

###############################################################################


def findPermutation2(pattern: str) -> List[int]:
    """寻找DI序列的字典序最小排列  O(nlogn)

    拓扑排序+pq
    沿着有向边走 数字逐渐递增 每次取出没用过的最小编号填入最小值
    """
    n = len(pattern) + 1
    deg, adjList = [0] * n, [[] for _ in range(n)]

    res = []
    for i, char in enumerate(pattern):
        if char == "I":
            deg[i + 1] += 1
            adjList[i].append(i + 1)
        else:
            deg[i] += 1
            adjList[i + 1].append(i)

    pq = [i for i in range(n) if deg[i] == 0]
    min_, res = 1, [0] * n
    while pq:
        cur = heappop(pq)  # 当前可用的最小序号
        res[cur] = min_
        min_ += 1
        for next in adjList[cur]:
            deg[next] -= 1
            if deg[next] == 0:
                heappush(pq, next)
    return res


print(findPermutation2("DI"))

###############################################################################


def smallestNumber(pattern: str) -> str:
    """
    寻找DI序列的字典序最小排列 len(s)<=8 可用数字为1-9 O(n!*n)
    9!=362880 全排列解法
    注意permutations是按字典序来的 找到就返回
    """
    n = len(pattern) + 1
    for cand in permutations(range(1, 10), n):
        for i, (pre, cur) in enumerate(pairwise(cand)):
            if pattern[i] == "D" and pre < cur:
                break
            elif pattern[i] == "I" and pre > cur:
                break
        else:
            return "".join(map(str, cand))
    return ""


print(smallestNumber("DI"))
