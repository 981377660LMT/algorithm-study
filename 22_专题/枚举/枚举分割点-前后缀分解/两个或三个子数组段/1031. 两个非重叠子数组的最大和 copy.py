from typing import List

INF = int(1e20)

# 给出非负整数数组 A ，返回两个非重叠（连续）子数组中元素的最大和，子数组的长度分别为 L 和 M
class Solution:
    def maxSumTwoNoOverlap(self, nums: List[int], left: int, right: int):
        def main(nums: List[int], len1: int, len2: int) -> int:
            """len1取前缀,len2取后缀"""
            preMax = [-INF] * n
            curSum = 0
            for i in range(n):
                curSum += nums[i]
                if i - (len1 - 1) >= 0:
                    preMax[i] = max(preMax[i - 1], curSum)
                    curSum -= nums[i - (len1 - 1)]

            nums = nums[::-1]
            suffixMax = [-INF] * n
            curSum = 0
            for i in range(n):
                curSum += nums[i]
                if i - (len2 - 1) >= 0:
                    suffixMax[i] = max(suffixMax[i - 1], curSum)
                    curSum -= nums[i - (len2 - 1)]
            suffixMax = suffixMax[::-1]

            return max(preMax[i] + suffixMax[i + 1] for i in range(n - 1))

        n = len(nums)
        return max(main(nums, left, right), main(nums, right, left))


# print(Solution().maxSumTwoNoOverlap([0, 6, 5, 2, 2, 5, 1, 9, 4], 1, 2))
# print(Solution().maxSumTwoNoOverlap([1, 0, 1], 1, 1))
# print(Solution().maxSumTwoNoOverlap([1, 0, 3], 1, 2))
# print(Solution().maxSumTwoNoOverlap([-2, 0, 2], 1, 2))
print(Solution().maxSumTwoNoOverlap([1, -2, 0, 0], 2, 1))
