#  * 786. 第 K 个最小的素数分数.py
from typing import List
from heapq import heappush, heappop

# 2 <= len(nums) <= 10000.
# 超时


class Solution:
    def smallestDistancePair(self, nums: List[int], k: int) -> int:
        def push(i: int, j: int):
            if 0 <= i and j < len(nums) and (i, j) not in visited:
                heappush(pq, (abs(nums[i] - nums[j]), i, j))
                visited.add((i, j))

        nums.sort()

        pq = []
        visited = set()
        for i in range(len(nums) - 1):
            heappush(pq, (abs(nums[i] - nums[i + 1]), i, i + 1))
            visited.add((i, i + 1))

        for _ in range(k - 1):
            val, i, j = heappop(pq)
            push(i - 1, j)
            push(i, j + 1)

        val, i, j = pq[0]
        return val


print(Solution().smallestDistancePair(nums=[1, 3, 1], k=1))
