#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# 最多有多少数字不小于x
# @param n int整型
# @param x int整型
# @param a int整型一维数组
# @return int整型
#
from typing import List


class Solution:
    def arrange(self, n, x, a: List[int]):
        # write code here
        a.sort(reverse=True)
        res = 0
        curSum = 0
        for i in range(n):
            curSum += a[i]
            if curSum // (i + 1) < x:
                break
            res += 1
        return res
