from collections import defaultdict

# 如果一个序列中 至少有三个元素 ，并且任意两个相邻元素之差相同，则称该序列为等差序列。
class Solution(object):
    def numberOfArithmeticSlices(self, nums):
        n = len(nums)
        dp = defaultdict(int)
        res = 0
        for i in range(1, n):
            for j in range(i):
                diff = nums[i] - nums[j]
                res += dp[(j, diff)]
                diff[(i, diff)] += dp[(j, diff)] + 1
        return res


print(Solution().numberOfArithmeticSlices([2, 4, 6, 8, 10]))
# 解释：所有的等差子序列为：
# [2,4,6]
# [4,6,8]
# [6,8,10]
# [2,4,6,8]
# [4,6,8,10]
# [2,4,6,8,10]
# [2,6,10]

