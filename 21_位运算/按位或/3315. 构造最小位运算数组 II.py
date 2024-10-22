# 3315. 构造最小位运算数组 II
# https://leetcode.cn/problems/construct-the-minimum-bitwise-array-ii/description/
# 给你一个长度为 n 的数组 nums 。你的任务是返回一个长度为 n 的数组 ans ，
# 对于每个下标 i ，以下 条件 均成立：
# ans[i] OR (ans[i] + 1) == nums[i]
# 除此以外，你需要 最小化 结果数组里每一个 ans[i] 。
# 如果没法找到符合 条件 的 ans[i] ，那么 ans[i] = -1 。
#
# x|(x+1) 的本质是把二进制最右边的 0 置为 1。


from typing import List


def lowbit(x: int) -> int:
    """二进制中最低位的1的位置.

    >>> [lowbit(x) for x in range(5)]
    [-1, 0, 1, 0, 2]
    """
    return -1 if not x else (x & -x).bit_length() - 1


class Solution:
    def minBitwiseArray(self, nums: List[int]) -> List[int]:
        def f(x: int) -> int:
            low0 = lowbit(~x)
            if low0 == 0:
                return -1
            x &= ~(1 << (low0 - 1))
            return x

        return [f(x) for x in nums]
