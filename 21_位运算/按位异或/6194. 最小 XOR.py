"""最小异或"""

# 给你两个正整数 num1 和 num2 ，找出满足下述条件的整数 x ：

# !x 的置位数和 num2 相同，且
# !x XOR num1 的值 最小
# 注意 XOR 是按位异或运算。

# 整数的 置位数 是其二进制表示中 1 的数目。


class Solution:
    def minimizeXor(self, num1: int, num2: int) -> int:
        count = num2.bit_count()
        res = 0
        for bit in range(32, -1, -1):  # !消除高位
            if (num1 >> bit) & 1:
                if count > 0:
                    count -= 1
                    res |= 1 << bit

        for bit in range(32):  # !补充低位
            if not count:
                break
            if not ((res >> bit) & 1):
                res |= 1 << bit
                count -= 1
        return res


print(Solution().minimizeXor(num1=3, num2=5))
