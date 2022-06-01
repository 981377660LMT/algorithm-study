from collections import defaultdict
from itertools import dropwhile, pairwise, takewhile, zip_longest
from typing import List


class Solution:
    def isAlienSorted(self, words: List[str], order: str) -> bool:
        rank = defaultdict(int, {v: i for i, v in enumerate(order)})
        for pre, cur in pairwise(words):
            for char1, char2 in zip(pre, cur):
                if char1 != char2:
                    if rank[char1] > rank[char2]:
                        return False
                    break
            else:
                if len(pre) > len(cur):
                    return False
        return True

    # def isAlienSorted2(self, words: List[str], order: str) -> bool:
    #     rank = defaultdict(int, {v: i + 1 for i, v in enumerate(order)})
    #     for w1, w2 in pairwise(words):
    #         # takewhile是取前缀 到某个条件不满足时停止取
    #         # dropwhile是取后缀 到某个条件不满足时开始取
    #         diffPair = next(
    #             dropwhile(lambda p: p[0] == p[1], zip_longest(w1, w2, fillvalue='')), None
    #         )
    #         # 取出第一个不同的字母对，如果第一个字母比第二个字母大，则返回False
    #         if diff and rank[diff[0]] > rank[diff[1]]:
    #             return False
    #     return True


print(Solution().isAlienSorted(["word", "world", "row"], "worldabcefghijkmnpqstuvxyz"))
# print(Solution().isAlienSorted2(["word", "world", "row"], "worldabcefghijkmnpqstuvxyz"))

