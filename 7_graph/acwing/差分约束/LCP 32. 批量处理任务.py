# LCP 32. 批量处理任务
# 2 <= tasks.length <= 10^5
# 0 <= tasks[i][0] <= tasks[i][1] <= 10^9
# 某实验室计算机待处理任务以 [start,end,period] 格式记于二维数组 tasks，
# 表示完成该任务的时间范围为起始时间 start 至结束时间 end 之间，需要计算机投入 period 的时长

# !处于开机状态的计算机可同时处理任意多个任务，请返回电脑最少开机多久，可处理完所有任务。
# https://leetcode-cn.com/problems/t3fKg1/solution/10xing-jie-jue-zhan-dou-by-foxtail-ke2e/


# 某实验室计算机待处理任务以 [start,end,period] 格式记于二维数组 tasks，
# 表示完成该任务的时间范围为起始时间 start 至结束时间 end 之间，
# 需要计算机投入 period 的时长
# 1. 离散化
# 2. 差分约束

# !python TLE 了

from collections import defaultdict, deque
from typing import List, Mapping, Tuple

INF = int(1e18)


class Solution:
    def processTasks(self, tasks: List[List[int]]) -> int:
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


print(Solution().processTasks([[1, 3, 2], [2, 5, 3], [5, 6, 2]]))

# 输出：4

# 解释：
# tasks[0] 选择时间点 2、3；
# tasks[1] 选择时间点 2、3、5；
# tasks[2] 选择时间点 5、6；
# 因此计算机仅需在时间点 2、3、5、6 四个时刻保持开机即可完成任务。
