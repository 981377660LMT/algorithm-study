# 6162. 收集垃圾的最少总时间
# 前缀和+哈希表记录每种字符的最后一个位置

from collections import defaultdict
from itertools import accumulate
from typing import List


class Solution:
    def garbageCollection(self, garbage: List[str], travel: List[int]) -> int:
        last = defaultdict(int)  # !每种字符的最后一个位置
        res = 0
        for i, word in enumerate(garbage):
            res += len(word)
            for char in word:
                last[char] = i

        preSum = [0] + list(accumulate(travel))
        for index in last.values():
            res += preSum[index]
        return res
