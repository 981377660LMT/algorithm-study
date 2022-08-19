from collections import defaultdict, deque
from typing import List, Sequence, Tuple

# LCP 32. 批量处理任务
# 2 <= tasks.length <= 10^5
# 0 <= tasks[i][0] <= tasks[i][1] <= 10^9
# 某实验室计算机待处理任务以 [start,end,period] 格式记于二维数组 tasks，
# 表示完成该任务的时间范围为起始时间 start 至结束时间 end 之间，需要计算机投入 period 的时长

# 处于开机状态的计算机可同时处理任意多个任务，请返回电脑最少开机多久，可处理完所有任务。
# https://leetcode-cn.com/problems/t3fKg1/solution/10xing-jie-jue-zhan-dou-by-foxtail-ke2e/


INF = 0x7FFFFFFF

# 某实验室计算机待处理任务以 [start,end,period] 格式记于二维数组 tasks，
# 表示完成该任务的时间范围为起始时间 start 至结束时间 end 之间，
# 需要计算机投入 period 的时长

# 1. 离散化
# 2. 差分约束


class Solution:
    def processTasks(self, tasks: List[List[int]]) -> int:
        def spfa(n: int, adjMap: Sequence) -> Tuple[bool, List[int]]:
            """spfa求单源最长路顺便判断正环"""
            dist = [0] * (n)
            queue = deque(list(range(n)))  # 这里很重要 每个点从0开始
            count = [0] * n
            isInqueue = [True] * n

            while queue:
                cur = queue.popleft()
                isInqueue[cur] = False

                for next, weight in adjMap[cur]:
                    if dist[cur] + weight > dist[next]:
                        dist[next] = dist[cur] + weight
                        count[next] = count[cur] + 1
                        if count[next] >= n + 1:
                            return False, []
                        if not isInqueue[next]:
                            isInqueue[next] = True
                            # 栈优化
                            queue.appendleft(next)

            return True, dist

        allNums = set()
        for start, end, period in tasks:
            allNums.add(start - 1)
            allNums.add(end)
        nums = sorted(allNums)

        n = len(allNums)
        mapping = {num: i for i, num in enumerate(nums)}

        adjList = [[] for _ in range(n)]
        for start, end, period in tasks:
            u = mapping[start - 1]
            v = mapping[end]
            adjList[u].append((v, period))

        for i in range(1, n):
            adjList[i - 1].append((i, 0))
            adjList[i].append((i - 1, nums[i - 1] - nums[i]))

        ok, dist = spfa(n, adjList)
        if ok:
            return dist[-1]
        return -1


print(Solution().processTasks([[1, 3, 2], [2, 5, 3], [5, 6, 2]]))

# 输出：4

# 解释：
# tasks[0] 选择时间点 2、3；
# tasks[1] 选择时间点 2、3、5；
# tasks[2] 选择时间点 5、6；
# 因此计算机仅需在时间点 2、3、5、6 四个时刻保持开机即可完成任务。
