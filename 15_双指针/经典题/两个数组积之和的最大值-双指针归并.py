# 0 ≤ n ≤ 100,000
# nums = [-5, 3, 2]
# multipliers = [-3, 1]
# We can match -5 with -3 and 3 with 1 to get -5 * -3 + 3 * 1.

# 不断选取两个数相乘 直到某个数组用尽
# 求和的最大值

from collections import deque


class Solution:
    def solve(self, nums, multipliers):
        nums = deque(sorted(nums))
        multipliers = deque(sorted(multipliers))
        res = 0
        while nums and multipliers:
            p1, p2 = nums[0] * multipliers[0], nums[-1] * multipliers[-1]
            if p1 > p2:
                res += p1
                nums.popleft()
                multipliers.popleft()
            else:
                res += p2
                nums.pop()
                multipliers.pop()
        return res
