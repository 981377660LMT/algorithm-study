from typing import List
from math import comb

# 二叉搜索树的方案数
MOD = 10 ** 9 + 7
# 此题与1916相似
# 子树1排序数*子树2排序数*组内保持顺序合并数组的方式comb


class Solution:
    def numOfWays(self, nums: List[int]) -> int:
        def countWays(nums: List[int]) -> int:
            if len(nums) <= 2:
                return 1
            left = [v for v in nums if v < nums[0]]
            right = [v for v in nums if v > nums[0]]
            return comb(len(left) + len(right), len(left)) * countWays(left) * countWays(right)

        # 最后记得要减1（去掉自身）。
        return (countWays(nums) - 1) % MOD


print(Solution().numOfWays([3, 4, 5, 1, 2]))

