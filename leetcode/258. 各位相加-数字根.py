# 258. 各位相加
# https://leetcode.cn/problems/add-digits/description/
#
# 极限脱出999里的数字根(数字根是指一个整数的各位数字反复相加，直到得到一个一位数)
# digital_root = 1 + (n - 1) % 9
#
# !对于 Python 来说，−1 mod 9 = 8，所以 Python 需要特判 num=0 的情况
# !X = 100a + 10b + c = 99a + 9b + (a+b+c)


class Solution:
    def addDigits(self, num: int) -> int:
        return 1 + (num - 1) % 9 if num else 0
