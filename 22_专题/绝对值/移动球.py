# 球用1表示
# 求出以每个位置为终点 移动球的花费


class Solution:
    def solve(self, nums):
        suffixSum = [0 for _ in range(len(nums))]
        prefixSum = [0 for _ in range(len(nums))]

        prefixCount = 0
        for i in range(len(nums)):
            prefixSum[i] = prefixCount + (prefixSum[i - 1] if i - 1 >= 0 else 0)
            if nums[i] == 1:
                prefixCount += 1

        suffixCount = 0
        for i in range(len(nums) - 1, -1, -1):
            suffixSum[i] = suffixCount + (suffixSum[i + 1] if len(nums) > i + 1 else 0)
            if nums[i] == 1:
                suffixCount += 1

        res = [prefixSum[i] + suffixSum[i] for i in range(len(nums))]

        return res


print(Solution().solve(nums=[1, 1, 0, 1]))
# [4, 3, 4, 5]
