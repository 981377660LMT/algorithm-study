# 100213. 按距离统计房屋对数目 II
# https://leetcode.cn/contest/weekly-contest-381/problems/count-the-number-of-houses-at-a-certain-distance-ii/
# 给你三个 正整数 n 、x 和 y 。
# 在城市中，存在编号从 1 到 n 的房屋，由 n 条街道相连。对所有 1 <= i < n ，都存在一条街道连接编号为 i 的房屋与编号为 i + 1 的房屋。另存在一条街道连接编号为 x 的房屋与编号为 y 的房屋。
# 对于每个 k（1 <= k <= n），你需要找出所有满足要求的 房屋对 [house1, house2] ，即从 house1 到 house2 需要经过的 最少 街道数为 k 。
# 返回一个下标从 1 开始且长度为 n 的数组 result ，其中 result[k] 表示所有满足要求的房屋对的数量，即从一个房屋到另一个房屋需要经过的 最少 街道数为 k 。
# 注意，x 与 y 可以 相等 。


from collections import deque
from itertools import accumulate
from typing import List

INF = int(1e18)


class DiffArray:
    """差分维护区间修改，区间查询."""

    __slots__ = ("_diff", "_dirty")

    def __init__(self, n: int) -> None:
        self._diff = [0] * (n + 1)
        self._dirty = False

    def add(self, start: int, end: int, delta: int) -> None:
        """区间 `[start,end)` 加上 `delta`."""
        if start < 0:
            start = 0
        if end >= len(self._diff):
            end = len(self._diff) - 1
        if start >= end:
            return
        self._dirty = True
        self._diff[start] += delta
        self._diff[end] -= delta

    def build(self) -> None:
        if self._dirty:
            self._diff = list(accumulate(self._diff))
            self._dirty = False

    def get(self, pos: int) -> int:
        """查询下标 `pos` 处的值."""
        self.build()
        return self._diff[pos]

    def getAll(self) -> List[int]:
        self.build()
        return self._diff[:-1]


def distPairOnCycle(n: int) -> List[int]:
    """环上两点距离为k的点对数."""
    res = [0] * (n // 2 + 1)
    for k in range(1, n // 2 + 1):
        res[k] = n
    if n % 2 == 0:
        res[n // 2] = n // 2
    return res


def bfs(start: int, adjList: List[List[int]]) -> List[int]:
    n = len(adjList)
    dist = [INF] * n
    dist[start] = 0
    queue = deque([start])
    while queue:
        cur = queue.popleft()
        for next in adjList[cur]:
            cand = dist[cur] + 1
            if cand < dist[next]:
                dist[next] = cand
                queue.append(next)
    return dist


class Solution:
    def countOfPairs(self, n: int, x: int, y: int) -> List[int]:
        if x > y:
            x, y = y, x
        x -= 1
        y -= 1

        # !左侧点[0,x)，环上点[x,y+1)， 右侧点[y+1,n)
        adjList = [[] for _ in range(n)]
        for i in range(n - 1):
            adjList[i].append(i + 1)
            adjList[i + 1].append(i)
        adjList[x].append(y)
        adjList[y].append(x)

        leftChainLength = x
        rightChainLength = n - y - 1
        cycleLength = y + 1 - x

        res = [0] * (n + 1)

        # leftChain、cycle、rightChain

        # !1.leftChain、rightChain 内部
        for d in range(1, leftChainLength):
            res[d] += leftChainLength - d
        for d in range(1, rightChainLength):
            res[d] += rightChainLength - d

        # !2.cycle 内部
        distPair = distPairOnCycle(cycleLength)
        for d, v in enumerate(distPair):
            res[d] += v

        bit = DiffArray(n)
        leftDist = bfs(x - 1, adjList) if x > 0 else []
        rightDist = bfs(y + 1, adjList) if y + 1 < n else []

        # !3.leftChain 到 cycle
        if leftChainLength > 0:
            counter = [0] * n
            for i in range(x, y + 1):
                counter[leftDist[i]] += 1
            for d, v in enumerate(counter):
                if v > 0:
                    bit.add(d, d + leftChainLength, v)

        # !4.rightChain 到 cycle
        if rightChainLength > 0:
            counter = [0] * n
            for i in range(x, y + 1):
                counter[rightDist[i]] += 1
            for d, v in enumerate(counter):
                if v > 0:
                    bit.add(d, d + rightChainLength, v)

        # !5.leftChain 到 rightChain
        if leftChainLength > 0 and rightChainLength > 0:
            counter = [0] * n
            for i in range(y + 1, n):
                counter[leftDist[i]] += 1
            for d, v in enumerate(counter):
                if v > 0:
                    bit.add(d, d + leftChainLength, v)

        for i in range(n):
            res[i] += bit.get(i)

        return [v * 2 for v in res[1:]]


if __name__ == "__main__":
    print(Solution().countOfPairs(4, 1, 2))
    # n = 5, x = 2, y = 4
    print(Solution().countOfPairs(5, 2, 4))
