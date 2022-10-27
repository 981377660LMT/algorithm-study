# 1 <= A.length <= 50000
# 滑动窗口最值问题
# 1.如果数组中的数据均为非负数的话，那么就对应常规的子数组和问题，可以使用滑动窗口来解决
# 209. 长度最小的子数组
# 但是添加了负数之后，窗口的滑动便丢失了单向性，因此无法使用滑动窗口解决。

from heapq import heappop, heappush
from typing import List
from collections import deque
from sortedcontainers import SortedList


INF = int(1e18)


class Solution:
    def shortestSubarray(self, nums: List[int], k: int) -> float:
        # 存储前缀和，单增；如果加入的前缀和减去队首的前缀和>=k 那么队首就找到了以他开始的最短的子数组，队首就可以退位了
        # 如果加入的前缀和小于队尾的前缀和 直接删除队尾 因为队尾找到符合题意的子数组还得比后面多带个负数 肯定不是最短的
        res = INF
        curSum = 0
        queue = deque([(-1, 0)])  # (index, preSum)
        for i, num in enumerate(nums):
            curSum += num
            while queue and curSum - queue[0][1] >= k:
                preI, _ = queue.popleft()
                res = min(res, i - preI)
            # 维护monoQueue单调性
            while queue and curSum <= queue[-1][1]:
                queue.pop()
            queue.append((i, curSum))  # type: ignore
        return res if res != INF else -1

    def shortestSubarray2(self, nums: List[int], k: int) -> float:
        """SortedList来维护有序性"""
        res = INF
        curSum = 0
        sl = SortedList([(0, -1)])
        for i, num in enumerate(nums):
            curSum += num
            while sl and curSum - sl[0][0] >= k:
                _, preI = sl.pop(0)
                res = min(res, i - preI)
            sl.add((curSum, i))
        return res if res != INF else -1

    def shortestSubarray3(self, nums: List[int], k: int) -> float:
        """pq来维护有序性"""
        res = INF
        curSum = 0
        pq = [(0, -1)]
        for i, num in enumerate(nums):
            curSum += num
            while pq and curSum - pq[0][0] >= k:
                _, preI = heappop(pq)
                res = min(res, i - preI)
            heappush(pq, (curSum, i))
        return res if res != INF else -1
