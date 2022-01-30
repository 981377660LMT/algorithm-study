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

# 平方是增长最快的方式，所以优先平方操作
# https://blog.nowcoder.net/n/78da997ed3aa46629885612e91a53f71?f=comment
class Solution:
    def solve(self, n, m) -> int:
        # write code here
        ...


print(Solution().solve(3, 10))
print(Solution().solve(1, 10))
print(Solution().solve(24, 500))

# 1.在当前数字的基础上加一，如：4转化为5
# 2.在当前数字的基础上减一，如：4转化为3
# 3.将当前数字变成它的平方，如：4转化为16

