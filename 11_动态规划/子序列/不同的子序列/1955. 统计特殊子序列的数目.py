from typing import List

# 特殊序列 是由 正整数 个 0 ，紧接着 正整数 个 1 ，最后 正整数 个 2 组成的序列。
# 回 不同特殊子序列的数目 。
# 1 <= nums.length <= 105
# https://leetcode-cn.com/problems/count-number-of-special-subsequences/solution/dong-tai-gui-hua-by-endlesscheng-4onu/

MOD = int(1e9 + 7)


class Solution:
    def countSpecialSubsequences(self, nums: List[int]) -> int:
        endswith = [0, 0, 0]
        for num in nums:
            # choose or not
            if num == 0:
                endswith[0] = endswith[0] + (endswith[0] + 1)
                endswith[0] %= MOD
            elif num == 1:
                endswith[1] = endswith[1] + (endswith[1] + endswith[0])
                endswith[1] %= MOD
            else:
                endswith[2] = endswith[2] + (endswith[2] + endswith[1])
                endswith[2] %= MOD

        return endswith[2]


print(Solution().countSpecialSubsequences(nums=[0, 1, 2, 2]))
# 输出：3
# 解释：特殊子序列为 [0,1,2,2]，[0,1,2,2] 和 [0,1,2,2] 。
