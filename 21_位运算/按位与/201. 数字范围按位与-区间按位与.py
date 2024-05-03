# find the LCP of the binary representations of left and right.
# rangeAdd 区间[left, right]的按位与


# !区间按位与


def rangeBitwiseAnd(left: int, right: int) -> int:
    if left == right:
        return left
    count = (left ^ right).bit_length()
    return (left >> count) << count


def rangeBitwiseAnd2(left: int, right: int) -> int:
    while left < right:
        right = right & (right - 1)
    return right


# 3125. 使得按位与结果为 0 的最大数字
# https://leetcode.cn/problems/maximum-number-that-makes-result-of-bitwise-and-zero/description/
# 给定一个整数 n，返回 最大的 整数 x 使得 x <= n，
# 并且所有在范围 [x, n] 内的数组的按位 AND 为 0。
class Solution:
    def maxNumber(self, n: int) -> int:
        return (1 << (n.bit_length() - 1)) - 1

        def check(mid: int) -> bool:
            return rangeBitwiseAnd(mid, n) == 0

        left, right = 0, n
        ok = False
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
                ok = True
            else:
                right = mid - 1
        return right
