from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def divisorSubstrings(self, num: int, k: int) -> int:
        """注意子字符串长度为 k，要算清楚"""
        res = 0
        n = len(str(num))
        for start in range(n - k + 1):
            cur = int(str(num)[start : start + k])
            if cur == 0:
                continue
            if num % cur == 0:
                res += 1
        return res
