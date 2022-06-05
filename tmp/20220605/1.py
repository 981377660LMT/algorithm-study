from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minMaxGame(self, nums: List[int]) -> int:
        while len(nums) > 1:
            nextQueue = []
            for i in range(0, len(nums), 2):
                if i % 4 == 0:
                    nextQueue.append(min(nums[i], nums[i + 1]))
                else:
                    nextQueue.append(max(nums[i], nums[i + 1]))
            nums = nextQueue
        return nums[0]
