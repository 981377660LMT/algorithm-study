from itertools import accumulate
from typing import Counter, List

MOD = int(1e9 + 7)

# 你可以任意排列 nums 中的数字，请你返回所有查询结果之和的最大值。
# 1 <= n <= 105
# 总结:重叠个数降序排列*nums降序排列


class Solution:
    def maxSumRangeQuery(self, nums: List[int], requests: List[List[int]]) -> int:
        n = len(nums)
        diff = [0] * (n + 1)
        for left, right in requests:
            diff[left] += 1
            diff[right + 1] -= 1
        diff = list(accumulate(diff))

        res = 0
        for value, count in zip(sorted(diff[:-1]), sorted(nums)):
            res += value * count
            res %= MOD
        return res


print(Solution().maxSumRangeQuery(nums=[1, 2, 3, 4, 5], requests=[[1, 3], [0, 1]]))

# 输出：19
# 解释：一个可行的 nums 排列为 [2,1,3,4,5]，并有如下结果：
# requests[0] -> nums[1] + nums[2] + nums[3] = 1 + 3 + 4 = 8
# requests[1] -> nums[0] + nums[1] = 2 + 1 = 3
# 总和为：8 + 3 = 11。
# 一个总和更大的排列为 [3,5,4,2,1]，并有如下结果：
# requests[0] -> nums[1] + nums[2] + nums[3] = 5 + 4 + 2 = 11
# requests[1] -> nums[0] + nums[1] = 3 + 5  = 8
# 总和为： 11 + 8 = 19，这个方案是所有排列中查询之和最大的结果。
