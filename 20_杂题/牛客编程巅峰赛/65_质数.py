#
# 返回两个区间内各取一个值相乘是p的倍数的个数
# @param a int整型 第一个区间的左边界
# @param b int整型 第一个区间的右边界
# @param c int整型 第二个区间的左边界
# @param d int整型 第二个区间的右边界
# @param p int整型 质数
# @return long长整型
#
from math import ceil


class Solution:
    def numbers(self, a: int, b: int, c: int, d: int, p: int) -> int:
        """
        有一个质数p，和两个区间[a,b]，[c,d]，
        分别在两个区间中取一个数x，y。求有多少对(x,y)使得x∗y是p的倍数
        """
        # 求出两个区间中的p的倍数 容斥原理
        x1 = ceil(a / p)
        x2 = b // p
        x = x2 - x1 + 1

        y1 = ceil(c / p)
        y2 = d // p
        y = y2 - y1 + 1

        return x * (d - c + 1) + y * (b - a + 1) - x * y
