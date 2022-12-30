# 给你两个整数，n 和 start 。
# 数组 nums 定义为：nums[i] = start + 2*i（下标从 0 开始）且 n == nums.length 。
# 请返回 nums 中所有元素按位异或（XOR）后得到的结果。

# 区间偶数的异或和
class Solution:
    def xorOperation(self, n: int, start: int) -> int:
        begin, end = start >> 1, n & start & 1
        res = preXor(begin - 1) ^ preXor(begin + n - 1)
        return (res << 1) | end


def preXor(upper: int) -> int:
    """[0, upper]内所有数的异或"""
    mod = upper % 4
    if mod == 0:
        return upper
    if mod == 1:
        return 1
    if mod == 2:
        return upper + 1
    return 0


assert Solution().xorOperation(5, 0) == 8
# 数组 nums 为 [0, 2, 4, 6, 8]，其中 (0 ^ 2 ^ 4 ^ 6 ^ 8) = 8 。
