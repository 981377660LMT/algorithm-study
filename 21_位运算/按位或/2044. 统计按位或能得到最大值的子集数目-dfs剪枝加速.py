from operator import or_
from functools import reduce, lru_cache
from typing import List


class Solution:
    # 2088ms
    def countMaxOrSubsets2(self, nums: List[int]) -> int:
        target = reduce(or_, nums)
        res = 0

        for state in range(1 << (len(nums))):
            orVal = 0
            for i in range(len(nums)):
                if state & (1 << i):
                    orVal |= nums[i]
            if orVal == target:
                res += 1

        return res

    # cache 提前结束
    # 48ms
    def countMaxOrSubsets(self, nums: List[int]) -> int:
        target = reduce(or_, nums)

        @lru_cache(None)
        def dfs(i, state):
            """Return number of subsets to get target."""
            # 如果当前 或的结果 已经和 数组整体或的结果 相同，则剩余右侧数字的含空全自己组合均为答案
            # 提前结束,起到了剪枝的效果
            if state == target:
                return 2 ** (len(nums) - i)
            if i == len(nums):
                return 0
            return dfs(i + 1, state | nums[i]) + dfs(i + 1, state)

        return dfs(0, 0)


print(Solution().countMaxOrSubsets(nums=[3, 2, 1, 5]))
# 输出：6
# 解释：子集按位或可能的最大值是 7 。有 6 个子集按位或可以得到 7 ：
# - [3,5]
# - [3,1,5]
# - [3,2,5]
# - [3,2,1,5]
# - [2,5]
# - [2,1,5]
