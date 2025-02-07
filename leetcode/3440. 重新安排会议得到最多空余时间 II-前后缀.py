# 给你一个整数 eventTime 表示一个活动的总时长，这个活动开始于 t = 0 ，结束于 t = eventTime 。
# 同时给你两个长度为 n 的整数数组 startTime 和 endTime 。它们表示这次活动中 n 个时间 没有重叠 的会议，其中第 i 个会议的时间为 [startTime[i], endTime[i]] 。
# 你可以重新安排 至多 一个会议，安排的规则是将会议时间平移，且保持原来的 会议时长 ，你的目的是移动会议后 最大化 相邻两个会议之间的 最长 连续空余时间。
# 请你返回重新安排会议以后，可以得到的 最大 空余时间。
# !注意：重新安排会议以后，会议之间的顺序可以发生改变。
#
# !枚举移动的会议
# 将会议平移到最左或最右，或者将会议移动到左边或右边的一个空隙.

from typing import List


class Solution:
    def maxFreeTime(self, eventTime: int, startTime: List[int], endTime: List[int]) -> int:
        n = len(startTime)
        leftMax = [0] * n  # leftMax[i] 表示第 i 个会议左边的最大空余时间
        rightMax = [0] * n  # rightMax[i] 表示第 i 个会议右边的最大空余时间
        leftMax[0] = startTime[0]
        for i in range(1, n):
            leftMax[i] = max(leftMax[i - 1], startTime[i] - endTime[i - 1])
        rightMax[n - 1] = eventTime - endTime[n - 1]
        for i in range(n - 2, -1, -1):
            rightMax[i] = max(rightMax[i + 1], startTime[i + 1] - endTime[i])

        res = 0
        for i in range(n):
            left = endTime[i - 1] if i > 0 else 0
            right = startTime[i + 1] if i < n - 1 else eventTime
            len_ = endTime[i] - startTime[i]
            res = max(res, right - left - len_)
            if i > 0 and leftMax[i - 1] >= len_:
                res = max(res, right - left)
            if i < n - 1 and rightMax[i + 1] >= len_:
                res = max(res, right - left)
        return res
