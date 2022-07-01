from functools import lru_cache
from typing import List


MAPPING = {
    1: 1,
    8: 8,
    6: 9,
    9: 6,
}  # 注意还有 {0:0}


class Solution:
    def findStrobogrammatic(self, n: int) -> List[str]:
        @lru_cache(None)
        def dfs(remain: int) -> List[str]:
            if remain < 0:
                return []
            if remain == 1:
                return ["0", "1", "8"]

            res: List[str] = []
            for nextRes in dfs(remain - 2):
                for head, tail in MAPPING.items():
                    res.append(str(head) + nextRes + str(tail))
                if remain != n:
                    res.append("0" + nextRes + "0")
            return res

        return dfs(n)
