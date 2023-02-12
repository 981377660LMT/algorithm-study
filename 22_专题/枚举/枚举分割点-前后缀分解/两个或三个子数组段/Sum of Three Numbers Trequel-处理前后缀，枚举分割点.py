# 找出和最大的三元组
# Return the maximum possible nums[i] + nums[j] + nums[k]
# 其中nums[i]<=nums[j]<=nums[k]

# 右侧前缀和维护最大，左侧sortedList维护，因为对每个分割点，要在左边找不超过他的最大元素


from itertools import accumulate
from sortedcontainers import SortedList


class Solution:
    def solve(self, nums):
        suffixMax = list(accumulate(nums[::-1], max))[::-1]
        pre = SortedList([nums[0]])

        n, res = len(nums), 0

        for midIndex in range(1, n - 1):
            mid = nums[midIndex]
            leftPos = pre.bisect_right(mid) - 1
            if leftPos < 0:
                pre.add(mid)
                continue

            left = pre[leftPos]
            right = suffixMax[midIndex + 1]
            if left <= mid <= right:
                res = max(res, mid + left + right)
            pre.add(mid)

        return res


print(Solution().solve(nums=[9, 1, 5, 3, 4]))
