# 选若干个物品，凑成容量target，求物品编号组成的最大值
from functools import lru_cache
from typing import List, Tuple


def compare(s1: str, s2: str) -> bool:
    """数字字符串比较大小"""
    return s1 > s2 if len(s1) == len(s2) else len(s1) > len(s2)


class Solution:
    def largestNumber(self, cost: List[int], target: int) -> str:
        @lru_cache(None)
        def dfs(remain: int) -> Tuple[bool, str]:
            if remain <= 0:
                return (True, '') if remain == 0 else (False, '')

            res = (False, '0')
            for select in MAPPING:
                cost = MAPPING[select]
                isOk, next = dfs(remain - cost)
                if not isOk:
                    continue
                cand = select + next
                # if int(cand) > int(res[1]):
                if compare(cand, res[1]):
                    res = (True, cand)

            return res

        MAPPING = {str(i + 1): num for i, num in enumerate(cost)}
        res = dfs(target)
        dfs.cache_clear()
        return res[1]


print(Solution().largestNumber([4, 3, 2, 5, 6, 7, 2, 5, 5], 9))

