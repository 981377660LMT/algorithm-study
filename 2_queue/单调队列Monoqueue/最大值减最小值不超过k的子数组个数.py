# 滑窗加MonoQueue 似乎可以
# right越往右滑 条件越苛刻  left收回来 条件就越好
from typing import List
from MonoQueue import MonoQueue


class Solution:
    def solve(self, nums: List[int], k: int) -> int:
        """
        求最大值减最小值不超过k的子数组个数
        """
        n = len(nums)
        queue = MonoQueue()
        res = 0
        for i in range(n):
            queue.append(nums[i])
            while queue and queue.max - queue.min > k:
                queue.popleft()
            res += len(queue)
        return res


print(Solution().solve(nums=[1, 3, 6], k=3))

