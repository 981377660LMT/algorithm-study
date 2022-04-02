# n ≤ 100,000
# 至多删除一个元素后 包含最大最小值的子数组 最多有多少个

# Dont remove anything
# Remove the maximum number only if the count of that is 1
# Remove the minimum number only if the count of that is 1
from typing import List


class Solution:
    def solve(self, nums: List[int]):
        def countSublistWithMaxAndMin(nums: List[int]) -> int:
            """包含最大最小值的子数组个数：按照子数组右端点进行分类"""
            res = 0
            max_, min_ = max(nums), min(nums)
            maxIndex, minIndex = -1, -1

            for i, num in enumerate(nums):
                if num == max_:
                    maxIndex = i
                if num == min_:
                    minIndex = i

                if maxIndex != -1 and minIndex != -1:
                    # 以i为右端点的子数组个数
                    res += min(maxIndex, minIndex) + 1

            return res

        n = len(nums)
        if n <= 1:
            return n

        res = countSublistWithMaxAndMin(nums)

        if nums.count(min(nums)) == 1:
            index = nums.index(min(nums))
            res = max(res, countSublistWithMaxAndMin(nums[:index] + nums[index + 1 :]))

        if nums.count(max(nums)) == 1:
            index = nums.index(max(nums))
            res = max(res, countSublistWithMaxAndMin(nums[:index] + nums[index + 1 :]))

        return res


# print(Solution().solve(nums=[2, 1, 5, 1, 3, 9]))
print(Solution().solve(nums=[0, 0]))

# If we remove 9 we'd get [2, 1, 5, 1, 3] and there's eight sublists where it contains both the max and the min:

# [1, 5]
# [5, 1]
# [1, 5, 1]
# [2, 1, 5]
# [5, 1, 3]
# [1, 5, 1, 3]
# [2, 1, 5, 1]
# [2, 1, 5, 1, 3]
