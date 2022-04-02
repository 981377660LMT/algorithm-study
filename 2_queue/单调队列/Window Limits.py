# 是否对于 abs(i - j) < window 的任何索引，abs(nums[i] - nums[j]) ≤ limit


from collections import deque


class Solution:
    def solve(self, nums, window, limit):
        minQueue = deque()  # <idx, val>
        maxQueue = deque()  # <idx, val>
        for i, num in enumerate(nums):
            while minQueue and minQueue[-1][1] >= num:
                minQueue.pop()
            while maxQueue and maxQueue[-1][1] <= num:
                maxQueue.pop()

            minQueue.append((i, num))
            maxQueue.append((i, num))

            if minQueue[0][0] <= i - window:
                minQueue.popleft()
            if maxQueue[0][0] <= i - window:
                maxQueue.popleft()

            if maxQueue[0][1] - minQueue[0][1] > limit:
                return False

        return True
