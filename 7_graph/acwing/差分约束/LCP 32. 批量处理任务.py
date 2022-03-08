from typing import List
from heapq import heappush, heappop

# 2 <= tasks.length <= 10^5
# 0 <= tasks[i][0] <= tasks[i][1] <= 10^9
# 某实验室计算机待处理任务以 [start,end,period] 格式记于二维数组 tasks，
# 表示完成该任务的时间范围为起始时间 start 至结束时间 end 之间，需要计算机投入 period 的时长

# 处于开机状态的计算机可同时处理任意多个任务，请返回电脑最少开机多久，可处理完所有任务。
# https://leetcode-cn.com/problems/t3fKg1/solution/10xing-jie-jue-zhan-dou-by-foxtail-ke2e/


INF = 0x7FFFFFFF


class Solution:
    def processTasks(self, tasks: List[List[int]]) -> int:
        ...


print(Solution().processTasks([[1, 3, 2], [2, 5, 3], [5, 6, 2]]))

# 输出：4

# 解释：
# tasks[0] 选择时间点 2、3；
# tasks[1] 选择时间点 2、3、5；
# tasks[2] 选择时间点 5、6；
# 因此计算机仅需在时间点 2、3、5、6 四个时刻保持开机即可完成任务。

