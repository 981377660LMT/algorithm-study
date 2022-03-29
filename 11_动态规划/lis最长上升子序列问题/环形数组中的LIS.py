from bisect import bisect_left

# n ≤ 1,000
# n^2logn


class Solution:
    def solve(self, nums):
        def LIS(left, right):
            res = [nums[left]]
            for i in range(left + 1, right + 1):
                index = bisect_left(res, nums[i])
                if index >= len(res):
                    res.append(nums[i])
                else:
                    res[index] = nums[i]
            return len(res)

        n = len(nums)
        nums = nums + nums
        return max(LIS(x, x + n) for x in range(n))
