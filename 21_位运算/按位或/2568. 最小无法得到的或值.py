# 给你一个下标从 0 开始的整数数组 nums 。

# 如果存在一些整数满足 0 <= index1 < index2 < ... < indexk < nums.length ，
# 得到 nums[index1] | nums[index2] | ... | nums[indexk] = x ，
# 那么我们说 x 是 可表达的 。换言之，如果一个整数能由 nums 的某个子序列的或运算得到，
# 那么它就是可表达的。

# 请你返回 nums 不可表达的 最小非零整数 。

# 考虑二进制表示
# 不可表达的最小整数肯定是二的幂
# 反证:假设不是二的幂，那么肯定缺失某一个位上的1,也就是二的幂,矛盾
# 因此只需要找到第一个不存在的二进制幂次
from typing import List


class Solution:
    def minImpossibleOR(self, nums: List[int]) -> int:
        S = set(nums)
        for i in range(32):
            if 1 << i not in S:
                return 1 << i
        raise Exception("no answer")

    def minImpossibleOR2(self, nums: List[int]) -> int:
        """先处理出所有二的幂次"""
        mask = 0
        for num in nums:
            if num & (num - 1) == 0:
                mask |= num

        # lowbit of zero
        mask = ~mask  # 取反后找最低位的0
        return mask & -mask


assert Solution().minImpossibleOR([1, 2, 4]) == 8
assert Solution().minImpossibleOR2([1, 2, 4]) == 8
