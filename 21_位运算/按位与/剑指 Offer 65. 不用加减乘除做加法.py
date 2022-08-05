# 不用加减乘除实现加法
# !注意公式 (a & b) << 1 + (a ^ b) = a + b
# 两数之和=两数与*2+两数异或
class Solution:
    def add(self, a: int, b: int) -> int:
        return (a & b) * 2 + (a ^ b)
