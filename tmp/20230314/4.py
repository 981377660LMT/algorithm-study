from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 你有一台电脑，它可以 同时 运行无数个任务。给你一个二维整数数组 tasks ，其中 tasks[i] = [starti, endi, durationi] 表示第 i 个任务需要在 闭区间 时间段 [starti, endi] 内运行 durationi 个整数时间点（但不需要连续）。

# 当电脑需要运行任务时，你可以打开电脑，如果空闲时，你可以将电脑关闭。

# 请你返回完成所有任务的情况下，电脑最少需要运行多少秒。

from collections import defaultdict, deque
from typing import List, Mapping, Tuple

INF = int(1e18)


class Solution:
    def findMinimumTime(self, tasks: List[List[int]]) -> int:
        def spfa(n: int, adjMap: Mapping[int, Mapping[int, int]]) -> Tuple[bool, List[int]]:
            """spfa求虚拟节点为起点的单源最长路 并检测正环"""
            dist = [0] * n
            queue = deque(list(range(n)))
            count = [0] * n
            inQueue = [True] * n

            while queue:
                cur = queue.popleft()
                inQueue[cur] = False
                for next in adjMap[cur]:
                    weight = adjMap[cur][next]
                    cand = dist[cur] + weight
                    if cand > dist[next]:
                        dist[next] = cand
                        count[next] = count[cur] + 1
                        if count[next] >= n:
                            return False, []
                        if not inQueue[next]:
                            inQueue[next] = True
                            queue.appendleft(next)  # !栈优化

            return True, dist

        allNums = set()
        for start, end, period in tasks:
            allNums.add(start - 1)
            allNums.add(end)
        nums = sorted(allNums)

        n = len(allNums)
        mp = {num: i for i, num in enumerate(nums)}

        adjMap = defaultdict(lambda: defaultdict(lambda: -INF))
        for start, end, period in tasks:
            u = mp[start - 1]
            v = mp[end]
            adjMap[u][v] = max(adjMap[u][v], period)  # v - u >= period

        for i in range(1, n):
            adjMap[i - 1][i] = max(adjMap[i - 1][i], 0)  # i - (i-1) >= 0
            adjMap[i][i - 1] = max(
                adjMap[i][i - 1], nums[i - 1] - nums[i]
            )  # i - (i-1) <= nums[i] - nums[i-1]

        ok, dist = spfa(n, adjMap)
        if not ok:
            return -1
        return dist[n - 1]
