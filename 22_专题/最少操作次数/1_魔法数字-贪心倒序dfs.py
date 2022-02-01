#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
#
# @param n int整型 表示牛牛的数字
# @param m int整型 表示牛妹的数字
# @return int整型 最少的操作将 n 转为 m
# (1≤n,m≤1000) 数据量可以bfs
#
from collections import deque
from math import floor, sqrt


# 到达离m最近的平方数
# https://blog.nowcoder.net/n/78da997ed3aa46629885612e91a53f71?f=comment
class Solution:
    def solve(self, n, m) -> int:
        # write code here
        if n >= m:
            return n - m
        cand1 = floor(sqrt(m))
        cand2 = cand1 + 1
        if abs(m - cand1 ** 2) > abs(m - cand2 ** 2):
            cand = cand2
        else:
            cand = cand1
        return min(m - n, 1 + self.solve(n, cand) + abs(m - cand ** 2))


print(Solution().solve(3, 10))
print(Solution().solve(1, 10))
print(Solution().solve(24, 500))

# 1.在当前数字的基础上加一，如：4转化为5
# 2.在当前数字的基础上减一，如：4转化为3
# 3.将当前数字变成它的平方，如：4转化为16

