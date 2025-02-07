# 3439. 重新安排会议得到最多空余时间 I
# https://leetcode.cn/problems/reschedule-meetings-for-maximum-free-time-i/description/
# 给你一个整数 eventTime 表示一个活动的总时长，这个活动开始于 t = 0 ，结束于 t = eventTime 。
# 同时给你两个长度为 n 的整数数组 startTime 和 endTime 。它们表示这次活动中 n 个时间 没有重叠 的会议，其中第 i 个会议的时间为 [startTime[i], endTime[i]] 。
# 你可以重新安排 至多 k 个会议，安排的规则是将会议时间平移，且保持原来的 会议时长 ，你的目的是移动会议后 最大化 相邻两个会议之间的 最长 连续空余时间。
# 移动前后所有会议之间的 相对 顺序需要保持不变，而且会议时间也需要保持互不重叠。
# 请你返回重新安排会议以后，可以得到的 最大 空余时间。
#
# 相对顺序需要保持不变，这意味着我们只能合并相邻的空余时间段，所以重新安排至多 k 个会议等价于如下问题：
# !给你 n+1 个空余时间段，合并其中 k+1 个连续的空余时间段，得到的最大长度是多少？


from typing import List


class Solution:
    def maxFreeTime(self, eventTime: int, k: int, startTime: List[int], endTime: List[int]) -> int:
        def f(i: int) -> int:
            if i == 0:
                return startTime[i]
            if i == n:
                return eventTime - endTime[i - 1]
            return startTime[i] - endTime[i - 1]

        k += 1
        n = len(startTime)
        res, curSum = 0, 0
        for right in range(n + 1):
            curSum += f(right)
            if right >= k:
                curSum -= f(right - k)
            if right >= k - 1:
                res = max(res, curSum)
        return res
