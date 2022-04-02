# Find a pair i ≤ j that maximizes values[i] + values[j] + nums[j] - nums[i] and return the value.


class Solution:
    def solve(self, nums, values):
        res = preMax = float("-inf")
        for num, val in zip(nums, values):
            preMax = max(preMax, val - num)
            res = max(res, preMax + val + num)

        return res


print(Solution().solve(nums=[1, 20, 2, 11], values=[1, 2, 3, 4]))
