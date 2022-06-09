from collections import defaultdict, deque
from functools import lru_cache
from typing import List

# 单调栈求出每个元素作为严格最大值/非严格最小值的影响范围
# 建图，dijkstra 求 0 到 n-1 的最短路
# !注意到原图是一个拓扑图(无环)，因此求最短路也可以用 dp O(n)求出 但是不一定比dijkstra快


class Solution:
    def minCost(self, nums: List[int], costs: List[int]) -> int:
        """dijk求最短路	1188 ms"""
        n = len(nums)
        maxRange = getRange(nums, isMax=True, isRightStrict=True)
        minRange = getRange(nums, isMax=False, isRightStrict=False)
        adjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))
        for cur, (next1, next2) in enumerate(zip(maxRange, minRange)):
            next1, next2 = next1 + 1, next2 + 1
            adjMap[cur][next1] = costs[next1] if next1 < n else int(1e20)
            adjMap[cur][next2] = costs[next2] if next2 < n else int(1e20)
        return dijkstra(adjMap, 0, n - 1)

    def minCost2(self, nums: List[int], costs: List[int]) -> int:
        """dfs求最短路	1584 ms"""

        @lru_cache(None)
        def dfs(cur: int) -> int:
            if cur >= n - 1:
                return 0 if cur == n - 1 else int(1e20)
            return min((dfs(next) + adjMap[cur][next] for next in adjMap[cur]), default=int(1e20))

        n = len(nums)
        maxRange = getRange(nums, isMax=True, isRightStrict=True)
        minRange = getRange(nums, isMax=False, isRightStrict=False)
        adjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))
        for cur, (next1, next2) in enumerate(zip(maxRange, minRange)):
            next1, next2 = next1 + 1, next2 + 1
            adjMap[cur][next1] = costs[next1] if next1 < n else int(1e20)
            adjMap[cur][next2] = costs[next2] if next2 < n else int(1e20)
        return dfs(0)

    def minCost3(self, nums: List[int], costs: List[int]) -> int:
        """拓扑排序求最短路	1188 ms"""

        n = len(nums)
        maxRange = getRange(nums, isMax=True, isRightStrict=True)
        minRange = getRange(nums, isMax=False, isRightStrict=False)
        adjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))
        deg = defaultdict(int)
        for cur, (next1, next2) in enumerate(zip(maxRange, minRange)):
            next1, next2 = next1 + 1, next2 + 1
            adjMap[cur][next1] = costs[next1] if next1 < n else int(1e20)
            adjMap[cur][next2] = costs[next2] if next2 < n else int(1e20)
            deg[next1] += 1
            deg[next2] += 1

        queue = deque([0])
        dist = defaultdict(lambda: int(1e20), {0: 0})
        while queue:
            cur = queue.popleft()
            for next in adjMap[cur]:
                dist[next] = min(dist[next], dist[cur] + adjMap[cur][next])
                deg[next] -= 1
                if deg[next] == 0:
                    queue.append(next)

        return dist[n - 1]


def getRange(nums: List[int], *, isMax=False, isRightStrict=False,) -> List[int]:
    def compareRight(stackValue: int, curValue: int) -> bool:
        if isRightStrict and isMax:
            return stackValue <= curValue
        elif isRightStrict and not isMax:
            return stackValue >= curValue
        elif not isRightStrict and isMax:
            return stackValue < curValue
        else:
            return stackValue > curValue

    n = len(nums)
    rightMost = [n - 1] * n
    stack = []
    for i in range(n):
        while stack and compareRight(nums[stack[-1]], nums[i]):
            rightMost[stack.pop()] = i - 1
        stack.append(i)
    return rightMost


from collections import defaultdict
from heapq import heappop, heappush
from typing import DefaultDict, Hashable, List, Optional, TypeVar, overload

INF = int(1e20)
Vertex = TypeVar('Vertex', bound=Hashable)
Graph = DefaultDict[Vertex, DefaultDict[Vertex, int]]


@overload
def dijkstra(adjMap: Graph, start: Vertex) -> DefaultDict[Vertex, int]:
    ...


@overload
def dijkstra(adjMap: Graph, start: Vertex, end: Vertex) -> int:
    ...


def dijkstra(adjMap: Graph, start: Vertex, end: Optional[Vertex] = None):
    dist = defaultdict(lambda: INF, {start: 0})
    pq = [(0, start)]

    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        if end is not None and cur == end:
            return curDist
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + adjMap[cur][next]:
                dist[next] = dist[cur] + adjMap[cur][next]
                heappush(pq, (dist[next], next))

    return INF if end is not None else dist


print(Solution().minCost(nums=[0, 1, 2], costs=[1, 1, 1]))
