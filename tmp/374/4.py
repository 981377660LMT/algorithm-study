from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个整数 n 和一个下标从 0 开始的整数数组 sick ，数组按 升序 排序。

# 有 n 位小朋友站成一排，按顺序编号为 0 到 n - 1 。
# 数组 sick 包含一开始得了感冒的小朋友的位置。
# 如果位置为 i 的小朋友得了感冒，他会传染给下标为 i - 1 或者 i + 1 的小朋友，前提 是被传染的小朋友存在且还没有得感冒。每一秒中， 至多一位 还没感冒的小朋友会被传染。

# 经过有限的秒数后，队列中所有小朋友都会感冒。感冒序列 指的是 所有 一开始没有感冒的小朋友最后得感冒的顺序序列。请你返回所有感冒序列的数目。

# 由于答案可能很大，请你将答案对 109 + 7 取余后返回。

# 注意，感冒序列 不 包含一开始就得了感冒的小朋友的下标


# 每段空白区间无关，若不在两端答案为 $2 ^ {len - 1}$
class Solution:
    def numberOfSequence(self, n: int, sick: List[int]) -> int:
        if len(sick) == 0:
            return 0

        adjList = [[] for _ in range(n)]
        for i in range(n):
            if i > 0:
                adjList[i].append(i - 1)
            if i < n - 1:
                adjList[i].append(i + 1)

        visited = [False] * n
        for i in sick:
            visited[i] = True

        queue = deque(sick)
        topoCount = 1
        while queue:
            len_ = len(queue)
            nextCount = 0
            for v in queue:
                for next in adjList[v]:
                    if not visited[next]:
                        nextCount += 1
            if nextCount > 0:
                topoCount = (topoCount * nextCount) % MOD
            for _ in range(len_):
                cur = queue.popleft()
                for next in adjList[cur]:
                    if not visited[next]:
                        visited[next] = True
                        queue.append(next)
                        nextCount += 1

        return topoCount % MOD


# n = 5, sick = [0,4]

print(Solution().numberOfSequence(5, [0, 4]))
