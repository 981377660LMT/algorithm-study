from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def digitSum(self, s: str, k: int) -> str:
        while len(s) > k:
            groups = [s[i : i + k] for i in range(0, len(s), k)]
            chars = []
            for g in groups:
                sum_ = sum(int(c) for c in g)
                chars.append(str(sum_))
            s = ''.join(chars)
        return s

