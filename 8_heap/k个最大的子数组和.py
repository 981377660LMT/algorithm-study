# n<=1000 k<=100

# O(n^2*logk)
import heapq
from itertools import accumulate


class Solution:
    def solve(self, nums, k):
        preSum = [0] + list(accumulate(nums))
        heap = []
        for right in range(len(preSum)):
            for left in range(right):
                if len(heap) >= k:
                    heapq.heappushpop(heap, preSum[right] - preSum[left])
                else:
                    heapq.heappush(heap, preSum[right] - preSum[left])

        return sorted(heap)

