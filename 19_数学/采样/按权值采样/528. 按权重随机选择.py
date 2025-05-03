# 528. 按权重随机选择
# https://leetcode.cn/problems/random-pick-with-weight/description/
# 彩票调度

import random
from typing import List
from typing import List
from bisect import bisect_left
from itertools import accumulate

from AliasMethod import AliasMethod


class Solution:
    def __init__(self, w: List[int]):
        sum_ = sum(w)
        probs = [p / sum_ for p in w]
        self.alias = AliasMethod(probs)

    def pickIndex(self) -> int:
        return self.alias.pick()


class Solution2:
    def __init__(self, w: List[int]):
        self.pre = list(accumulate(w))
        self.total = sum(w)

    def pickIndex(self) -> int:
        x = random.randint(1, self.total)
        return bisect_left(self.pre, x)
