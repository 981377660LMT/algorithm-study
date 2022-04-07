# 使用给定的字符串，能组成多少个不同的回文
from collections import Counter
from math import factorial


class Solution:
    def solve(self, s):
        counter = Counter(s)
        if sum(v % 2 for v in counter.values()) > 1:
            return 0

        n = 0
        r = 1
        for v in counter.values():
            k = v // 2
            n += k
            r *= factorial(k)
        return (factorial(n) // r) % (10 ** 9 + 7)
