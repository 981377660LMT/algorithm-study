#
# 两个数表示答案
# @param n int整型 一次运输的冰激凌数量
# @param m int整型 总冰激凌数
# @param t int整型 一次运输的时间
# @param c int整型一维数组 表示每个冰激凌制作好时间<1e4
# @return int整型一维数组
#
from typing import List


# 贪心，直接先运送制作快的即可
class Solution:
    def icecream(self, n: int, m: int, t: int, c: List[int]):
        # write code here
        c.sort()
        div, mod = divmod(m, n)
        if mod == 0:
            count = div
            # 车上最后一个冰淇淋(制作时间最久)的编号
            carry = n - 1
        else:
            count = div + 1
            carry = mod - 1
        time = 0
        while carry < m:
            time = max(time, c[carry]) + 2 * t
            carry += n
        time -= t  # 减去最后一趟的时间
        return [time, count]

