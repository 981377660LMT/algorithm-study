# k次操作 每次操作可以为1个数加1
# 求k次操作后 最长连续子数组的长度

# n ≤ 100,000
# k < 2 ** 31

from collections import deque
from typing import List
from sortedcontainers import SortedList


class Solution:
    def solve1(self, nums: List[int], k: int):
        """sortedList TLE"""
        # 我们需要维护窗口内改变代价最大的那一个元素,可以用monoQueue,也可以sortedLst
        window = SortedList(key=lambda x: -x)
        res = 0

        left = 0
        curSum = 0
        for right, num in enumerate(nums):
            curSum += num
            window.add(num)

            # 统一为max的代价超出了,滑窗左移
            while window and (right - left + 1) * window[0] - curSum > k:
                curSum -= nums[left]
                window.discard(nums[left])
                left += 1

            res = max(res, right - left + 1)

        return res

    def solve(self, nums: List[int], k: int):
        """monoQueue AC"""
        # 我们需要维护窗口内改变代价最大的那一个元素,可以用monoQueue,也可以sortedLst
        window = deque()
        res = 0

        left = 0
        curSum = 0
        for right, num in enumerate(nums):
            curSum += num
            while window and nums[window[-1]] < num:
                window.pop()
            window.append(right)

            # 统一为max的代价超出了,滑窗左移
            while window and (right - left + 1) * nums[window[0]] - curSum > k:
                curSum -= nums[left]
                if window and window[0] == left:
                    window.popleft()
                left += 1

            res = max(res, right - left + 1)

        return res


print(Solution().solve(nums=[2, 4, 8, 5, 9, 6], k=6))
print(Solution().solve(nums=[2, 3, 1, 1], k=0))
