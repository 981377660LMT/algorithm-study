# 3630. 划分数组得到最大异或运算和与运算之和
# https://leetcode.cn/problems/partition-array-for-maximum-xor-and-and/description/
# 给你一个整数数组 nums。
# 将数组划分为 三 个（可以为空）子序列 A、B 和 C，使得 nums 中的每个元素 恰好 属于一个子序列。
# 你的目标是 最大化 以下值：XOR(A) + AND(B) + XOR(C)
# 1 <= nums.length <= 19
# 1 <= nums[i] <= 1e9
#
# !枚举 nums 的子集作为 B，这样我们只需专注 XOR 的问题

from typing import List
from 数组分成两部分最大化异或和 import solve


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maximizeXorAndXor(self, nums: List[int]) -> int:
        n = len(nums)
        res = 0
        for s in range(1 << n):
            and_, xorNums = -1, []
            for i, v in enumerate(nums):
                if (s >> i) & 1:
                    and_ &= v
                else:
                    xorNums.append(v)
            and_ = 0 if and_ == -1 else and_
            res = max2(res, solve(xorNums) + and_)
        return res


# [2,3,6,7]
print(Solution().maximizeXorAndXor([2, 3, 6, 7]))  # 8
