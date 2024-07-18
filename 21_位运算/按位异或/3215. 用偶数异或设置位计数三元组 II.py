# 3215. 用偶数异或设置位计数三元组 II
# https://leetcode.cn/problems/count-triplets-with-even-xor-set-bits-ii/description/
# !给定三个整数数组 a，b 和 c，返回组内元素按位异或二进制为1得位数有偶数个 的三元组 (a[i], b[j], c[k]) 的数量
#
# !异或不改变1位数量的奇偶性.
# !3个数异或结果有偶数个1，要么3个数都有偶数个1，要么其中1个数有偶数个1，另2个数都有奇数个1.

from itertools import product
from typing import List


class Solution:
    def tripletCount(self, a: List[int], b: List[int], c: List[int]) -> int:
        def cal(nums: List[int]) -> List[int]:
            res = [0, 0]
            for num in nums:
                res[num.bit_count() & 1] += 1
            return res

        counter1, counter2, counter3 = cal(a), cal(b), cal(c)
        res = 0
        for i, j, k in product(range(2), repeat=3):
            if not (i + j + k) & 1:
                res += counter1[i] * counter2[j] * counter3[k]
        return res
