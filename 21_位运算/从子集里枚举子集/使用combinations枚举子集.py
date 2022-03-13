import collections
from itertools import chain, combinations
from typing import Any, Collection, List


def powerset(collection: Collection[Any]):
    """求真子集"""
    return chain.from_iterable(combinations(collection, n) for n in range(len(collection)))


class Solution:
    def coopDevelop(self, skills: List[List[int]]) -> int:
        """使用示例"""
        counter = collections.Counter()
        n = len(skills)
        for skill in skills:
            counter[tuple(skill)] += 1
        res = n * (n - 1) // 2
        for count in counter.values():
            res -= count * (count - 1) // 2
        for skill in counter:
            for sub in powerset(skill):
                res -= counter[skill] * counter[sub]
        return res % int(1e9 + 7)

