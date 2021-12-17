# 剑指 Offer 14- II. 剪绳子 II
# n>4 的条件下尽量取3


# 你需要构造一个正整数 n ，它满足以下条件：
# n 质因数（质因数需要考虑重复的情况）的数目 不超过 primeFactors 个。
# n 好因子的数目最大化(n 的一个因子可以被 n 的每一个质因数整除)
# 比方说，如果 n = 12 ，那么它的质因数为 [2,2,3] ，那么 6 和 12 是好因子，但 3 和 4 不是。

# 1 <= primeFactors <= 109

# https://leetcode-cn.com/problems/maximize-number-of-nice-divisors/comments/861510


# 已知e1+e2+...+e_n为定值 (因子个数和)
# 按照乘法原理，num 所有好因子的个数即为e_1*e_2*e_3...e_n
# 由于n太大无法dp求解,所以只能用数学结论

M = int(1e9 + 7)


class Solution:
    def maxNiceDivisors(self, num: int) -> int:
        if num == 1:
            return 1
        if num == 2:
            return 2

        div, mod = divmod(num, 3)
        if mod == 0:
            return pow(3, div, M)
        if mod == 1:
            return pow(3, div - 1, M) * 4 % M
        if mod == 2:
            return pow(3, div, M) * 2 % M

