# rangeAdd 区间[left, right]的按位与


# !区间按位与
# find the LCP of the binary representations of left and right.
class Solution:
    def rangeBitwiseAnd(self, left: int, right: int) -> int:
        if left == right:
            return left
        count = (left ^ right).bit_length()
        return (left >> count) << count

    def rangeBitwiseAnd2(self, left: int, right: int) -> int:
        while left < right:
            right = right & (right - 1)
        return right
