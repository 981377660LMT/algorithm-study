# 1819. 序列中不同最大公约数的数目
# 统计数组中所有子序列的 gcd 的不同个数，复杂度 O(maxlog^2max)
# https://leetcode.cn/problems/number-of-different-subsequences-gcds/solution/ji-bai-100mei-ju-gcdxun-huan-you-hua-pyt-get7/
#
# 1<=nums.length<=1e5
# 1<=nums[i]<=2e5


from math import gcd
from typing import List


class Solution:
    def countDifferentSubsequenceGCDs(self, nums: List[int]) -> int:
        """计算并返回 nums 的所有 非空 子序列中 不同 最大公约数的 数目 。"""
        max_ = max(nums)
        has = [False] * (max_ + 1)
        for num in nums:
            has[num] = True

        # 枚举答案
        res = 0
        for i in range(1, max_ + 1):
            gcd_ = 0
            for j in range(i, max_ + 1, i):
                if has[j]:
                    gcd_ = gcd(gcd_, j)
                    if gcd_ == i:
                        res += 1
                        break

        return res


if __name__ == "__main__":
    print(Solution().countDifferentSubsequenceGCDs(nums=[6, 10, 3]))
    print(Solution().countDifferentSubsequenceGCDs(nums=[3, 6]))
