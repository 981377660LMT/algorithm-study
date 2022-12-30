# 267. 回文排列 II
# 全排列生成回文的一半
# 1 <= s.length <= 16
# s 仅由小写英文字母组成

from collections import Counter
from itertools import permutations
from typing import List


class Solution:
    def generatePalindromes(self, s: str) -> List[str]:
        counter = Counter(s)
        odd, even = [], []
        for k, v in counter.items():
            if v & 1:
                odd.append(k)
            even.extend([k] * (v // 2))

        if len(odd) > 1:
            return []

        res = set()
        for perm in permutations(even):
            res.add("".join(perm) + "".join(odd) + "".join(perm[::-1]))
        return list(res)


print(Solution().generatePalindromes("aabb"))
