# https://leetcode-cn.com/problems/t3fKg1/solution/10xing-jie-jue-zhan-dou-by-foxtail-ke2e/
# LCP 32. 批量处理任务
# 2 <= tasks.length <= 10^5
# 0 <= tasks[i][0] <= tasks[i][1] <= 10^9
# 某实验室计算机待处理任务以 [start,end,period] 格式记于二维数组 tasks，
# 表示完成该任务的时间范围为起始时间 start 至结束时间 end 之间，需要计算机投入 period 的时长
# !处于开机状态的计算机可同时处理任意多个任务，请返回电脑最少开机多久，可处理完所有任务。

from typing import List
from 差分约束 import DualShortestPath


class Solution:
    def processTasks(self, tasks: List[List[int]]) -> int:
        allNums = set()
        for s, e, p in tasks:
            allNums.add(s - 1)
            allNums.add(e)
        allNums = sorted(allNums)

        n = len(allNums)
        mp = {num: i for i, num in enumerate(allNums)}
        D = DualShortestPath(n + 10, min=True)
        for s, e, p in tasks:
            u, v = mp[s - 1], mp[e]  # v - u >= period
            D.addEdge(u, v, -p)
        for i in range(1, n):
            D.addEdge(i - 1, i, 0)  # Si>=Si-1
            D.addEdge(i, i - 1, allNums[i] - allNums[i - 1])  # Si-Si-1<=allNums[i]-allNums[i-1]

        res, ok = D.run()
        if not ok:
            return -1
        return res[n - 1]


print(Solution().processTasks([[1, 3, 2], [2, 5, 3], [5, 6, 2]]))
