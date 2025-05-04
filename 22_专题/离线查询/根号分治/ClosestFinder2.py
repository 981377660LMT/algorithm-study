# O(√qnlogn) 的二分查找+缓存查询结果解法
# https://leetcode.cn/problems/shortest-word-distance-ii/solutions/1941801/zheng-que-de-o-by-vclip-ypg6/
#
# 1. 记录每个字符串的出现的所有下标。查询时遍历一个字符串的所有下标，在另一个字符串的下标列表里二分查找最近的下标，用两个下标的差值更新最短距离。
# 2. 加了缓存的二分查找法的复杂度为O(√qnlogn).


from typing import List
from bisect import bisect_left


INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class WordDistance:
    __slots__ = ("_indexes", "_cache")

    def __init__(self, wordsDict: List[str]):
        indexes = dict()
        for i, w in enumerate(wordsDict):
            indexes.setdefault(w, []).append(i)
        self._indexes = indexes
        self._cache = dict()

    def shortest(self, word1: str, word2: str) -> int:
        hash_ = (word1, word2) if word1 < word2 else (word2, word1)
        res = self._cache.get(hash_, None)
        if res is not None:
            return res
        res = self._shortest(word1, word2)
        self._cache[hash_] = res
        return res

    def _shortest(self, word1: str, word2: str) -> int:
        pos1, pos2 = self._indexes.get(word1, []), self._indexes.get(word2, [])
        if len(pos1) > len(pos2):
            pos1, pos2 = pos2, pos1
        res = INF
        for p in pos1:
            i = bisect_left(pos2, p)
            if i < len(pos2):
                res = min2(res, pos2[i] - p)
            if i > 0:
                res = min2(res, p - pos2[i - 1])
        return res if res != INF else -1


# Your WordDistance object will be instantiated and called as such:
# obj = WordDistance(wordsDict)
# param_1 = obj.shortest(word1,word2)
