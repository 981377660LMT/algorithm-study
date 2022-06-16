# 使用给定的字符串，能组成多少个不同的回文
from collections import Counter
from math import factorial

# 所有字符组成的不同的回文串的个数
# return the number of distinct palindromes you can make using all characters.

# !找回文的一半有哪些
MOD = int(1e9 + 7)


class Solution:
    def solve(self, s: str):
        """回文总排列数除以各个排列数"""
        counter = Counter(s)
        if sum(v % 2 for v in counter.values()) > 1:
            return 0

        res = 0
        div = 1
        for count in counter.values():
            # 这个字符贡献的对数
            half = count // 2
            res += half
            div *= factorial(half)

        return (factorial(res) // div) % MOD


print(Solution().solve(s="abccb"))  # 2
# We can make "cbabc" and "bcacb"
