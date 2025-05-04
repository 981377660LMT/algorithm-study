# 2589. 完成所有任务的最少时间
# https://leetcode.cn/problems/minimum-time-to-complete-all-tasks/description/
# 你有一台电脑，它可以 同时 运行无数个任务。给你一个二维整数数组 tasks ，其中 tasks[i] = [starti, endi, durationi] 表示第 i 个任务需要在 闭区间 时间段 [starti, endi] 内运行 durationi 个整数时间点（但不需要连续）。
# 当电脑需要运行任务时，你可以打开电脑，如果空闲时，你可以将电脑关闭。
# 请你返回完成所有任务的情况下，电脑最少需要运行多少秒。
#
# !贪心的策略：
# !按任务结束时间排序后，后结束的任务应当先给其分配[start,end]范围内被先结束的任务占用的时间点，再给其按从end到start的顺序分配未被占用的时间点

from bisect import bisect_left
from typing import List


class Solution:
    def findMinimumTime(self, tasks: List[List[int]]) -> int:
        """
        O(nlogn)栈+二分.
        !将占用的时间点用不相交区间表示.
        O(logn) 查询后缀和.
        均摊O(1) 更新后缀.
        """
        tasks.sort(key=lambda x: x[1])
        stack = []
        presum = [0]
        for start, end, du in tasks:
            pos = bisect_left(stack, (start,))
            du -= presum[-1] - presum[pos]
            if pos > 0 and start <= (r := stack[pos - 1][1]):
                du -= r - start + 1
            if du <= 0:
                continue
            while stack and end - stack[-1][1] < du:
                presum.pop()
                l, r = stack.pop()
                du += r - l + 1
            stack.append((end - du + 1, end))
            presum.append(presum[-1] + du)
        return presum[-1]

    def findMinimumTime2(self, tasks: List[List[int]]) -> int:
        """O(n*U)暴力."""
        tasks.sort(key=lambda x: x[1])
        run = [False] * (tasks[-1][1] + 1)
        for start, end, du in tasks:
            du -= sum(run[start : end + 1])  # !瓶颈1
            if du <= 0:
                continue  # finished
            for i in range(end, start - 1, -1):  # !瓶颈2
                if run[i]:
                    continue
                run[i] = True
                du -= 1
                if du == 0:
                    break
        return sum(run)
