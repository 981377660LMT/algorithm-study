# 3296. 移山所需的最少秒数-二分套二分
# https://leetcode.cn/problems/minimum-number-of-seconds-to-make-mountain-height-zero/description/
# 给你一个整数 mountainHeight 表示山的高度。
# 同时给你一个整数数组 workerTimes，表示工人们的工作时间（单位：秒）。
# 工人们需要 同时 进行工作以 降低 山的高度。对于工人 i :
# 山的高度降低 x，需要花费 workerTimes[i] + workerTimes[i] * 2 + ... + workerTimes[i] * x 秒。例如：
# 山的高度降低 1，需要 workerTimes[i] 秒。
# 山的高度降低 2，需要 workerTimes[i] + workerTimes[i] * 2 秒，依此类推。
# 返回一个整数，表示工人们使山的高度降低到 0 所需的 最少 秒数。
#
# !堆模拟，比较函数为工作时间之和.


from heapq import heapify, heappop, heappush
from typing import List


class Solution:
    def minNumberOfSeconds(self, mountainHeight: int, workerTimes: List[int]) -> int:
        n = len(workerTimes)
        times = [0] * n
        pq = [(v, v, i) for i, v in enumerate(workerTimes)]
        heapify(pq)
        for _ in range(mountainHeight):
            max_, v, i = heappop(pq)
            times[i] = max_
            v += workerTimes[i]
            max_ += v
            heappush(pq, (max_, v, i))
        return max(times)
