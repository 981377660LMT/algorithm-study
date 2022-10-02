from math import gcd


# !给你两个正整数 a 和 b ，返回 a 和 b 的 公因子的数目。
# 如果 x 可以同时整除 a 和 b ，则认为 x 是 a 和 b 的一个 公因子 。


class Solution:
    def commonFactors(self, a: int, b: int) -> int:
        gcd_ = gcd(a, b)  # !上界为最大公约数
        res = 0
        for i in range(1, gcd_ + 1):
            if a % i == 0 and b % i == 0:
                res += 1
        return res
