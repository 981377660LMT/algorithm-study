from typing import List
from heapq import heappush, heappop

# 2 <= tasks.length <= 10^5
# 0 <= tasks[i][0] <= tasks[i][1] <= 10^9
# 某实验室计算机待处理任务以 [start,end,period] 格式记于二维数组 tasks，
# 表示完成该任务的时间范围为起始时间 start 至结束时间 end 之间，需要计算机投入 period 的时长

# 处于开机状态的计算机可同时处理任意多个任务，请返回电脑最少开机多久，可处理完所有任务。
# https://leetcode-cn.com/problems/t3fKg1/solution/10xing-jie-jue-zhan-dou-by-foxtail-ke2e/

# 1. 最少开机：计算机很懒，不必须开机就不开机
# 2. 必须开机 => 当前是任务最迟开始时间，这时就开机
# 3. 对于一堆都已经开始的任务，当系统运行a时长后，所有任务的最晚启动时刻也会集体后移a（
# 4. 入队时，最早开始时间-res(为了抵消后面+res)，在队列中+res
INF = 0x7FFFFFFF


class Solution:
    def processTasks(self, tasks: List[List[int]]) -> int:
        tasks.append([INF, INF + 1, 1])
        # 化成左闭右开区间，便于处理最早开始时间
        tasks = sorted([s, e + 1, p] for s, e, p in tasks)

        res = 0
        pq = []
        for [start, end, cost] in tasks:
            # 需要处理任务了
            while pq and pq[0][0] + res < start:
                # 过期任务
                if pq[0][0] + res >= pq[0][1]:
                    heappop(pq)
                # 执行一个任务
                else:
                    res += min(pq[0][1], start) - (pq[0][0] + res)
            # 入队
            heappush(pq, (end - cost - res, end))

        return res


print(Solution().processTasks([[1, 3, 2], [2, 5, 3], [5, 6, 2]]))

# 输出：4

# 解释：
# tasks[0] 选择时间点 2、3；
# tasks[1] 选择时间点 2、3、5；
# tasks[2] 选择时间点 5、6；
# 因此计算机仅需在时间点 2、3、5、6 四个时刻保持开机即可完成任务。

