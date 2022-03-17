from operator import or_
from functools import reduce, lru_cache
from typing import List


class Solution:
    def countMaxOrSubsets2(self, nums: List[int]) -> int:
        target = reduce(or_, nums)
        xors = [0] * (1 << len(nums))
        res = 0
        # 时间复杂度为 O(1+2+4+...+2^(n-1)) = O(2^n)
        for i, num in enumerate(nums):
            for preState in range(1 << i):
                cur = xors[preState] | num
                xors[preState | (1 << i)] = cur
                if cur == target:
                    res += 1

        return res

    # cache 提前结束
    # 48ms
    def countMaxOrSubsets(self, nums: List[int]) -> int:
        @lru_cache(None)
        def dfs(index: int, state: int) -> int:
            """Return number of subsets to get target."""
            # 如果当前 或的结果 已经和 数组整体或的结果 相同，则剩余右侧数字的含空全自己组合均为答案
            # 提前结束,起到了剪枝的效果
            if state == target:
                return 2 ** (len(nums) - index)
            if index == len(nums):
                return 0
            return dfs(index + 1, state | nums[index]) + dfs(index + 1, state)

        target = reduce(or_, nums)
        return dfs(0, 0)


print(Solution().countMaxOrSubsets2(nums=[3, 2, 1, 5]))
# 输出：6
# 解释：子集按位或可能的最大值是 7 。有 6 个子集按位或可以得到 7 ：
# - [3,5]
# - [3,1,5]
# - [3,2,5]
# - [3,2,1,5]
# - [2,5]
# - [2,1,5]
