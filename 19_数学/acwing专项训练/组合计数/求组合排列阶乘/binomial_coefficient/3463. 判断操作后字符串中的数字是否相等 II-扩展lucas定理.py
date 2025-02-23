# 3463. 判断操作后字符串中的数字是否相等 II
# https://leetcode.cn/problems/check-if-digits-are-equal-in-string-after-operations-ii/description/
# 给你一个由数字组成的字符串 s 。重复执行以下操作，直到字符串恰好包含 两个 数字：
#
# !从第一个数字开始，对于 s 中的每一对连续数字，计算这两个数字的和 模 10。
# 用计算得到的新数字依次替换 s 的每一个字符，并保持原本的顺序。
# 如果 s 最后剩下的两个数字相同，则返回 true 。否则，返回 false。
#
# 杨辉三角/二项式系数
# !下标为i的数在第n层结果中的权重为C(n-1,i)

from BinomialCoefficient import BinomialCoefficient

MOD = 10
C = BinomialCoefficient(MOD)


class Solution:
    def hasSameDigits(self, s: str) -> bool:
        n = len(s)
        nums = [int(x) for x in s]
        sum1 = sum(nums[i] * C(n - 2, i) for i in range(n - 1)) % MOD
        sum2 = sum(nums[i + 1] * C(n - 2, i) for i in range(n - 1)) % MOD
        return sum1 == sum2
