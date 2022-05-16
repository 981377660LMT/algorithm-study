from itertools import chain, combinations
from typing import Collection, List

# n,m<=20
# 选择word使其无相同字符，求最长长度


def powerset(collection: Collection, isAll=True):
    """求(真)子集,时间复杂度O(n*2^n)

    默认求所有子集
    """
    upper = len(collection) + 1 if isAll else len(collection)
    return chain.from_iterable(combinations(collection, n) for n in range(upper))


class Solution:
    def maxLength(self, arr: List[str]) -> int:
        res = 0
        for p in powerset(arr):
            allChars = ''.join(w for w in p)
            if len(allChars) == len(set(allChars)):
                res = max(res, len(allChars))
        return res
