# 每次选中一个偶数，然后把这些数中与该数相等的数都除以2，例如现在有一个数组为[2,2,3][2,2,3]，那么牛牛可以执行一次操作，使得这个数组变为[1,1,3][1,1,3]。
# 牛牛现在想知道，对于任意的n个数，他最少需要操作多少次，使得这些数都变成奇数？
#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# 返回一个数，代表让这些数都变成奇数的最少的操作次数
# @param n int整型 代表一共有多少数
# @param a int整型一维数组 代表n个数字的值
# @return int整型
#
from typing import List


class Solution:
    def oddNumber(self, n: int, a: List[int]):
        # write code here
        ...
        res = 0
        visited = set()
        for num in a:
            while num % 2 == 0 and num not in visited:
                visited.add(num)
                num //= 2
                res += 1
        return res

