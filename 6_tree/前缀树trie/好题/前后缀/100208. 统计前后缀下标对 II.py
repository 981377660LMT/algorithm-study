# 3045. 统计前后缀下标对 I
# https://leetcode.cn/problems/count-prefix-and-suffix-pairs-ii/
# 当 str1 同时是 str2 的前缀和后缀时，isPrefixAndSuffix(str1, str2) 返回 true，否则返回 false。
# 以整数形式，返回满足 i < j 且 isPrefixAndSuffix(words[i], words[j]) 为 true 的下标对 (i, j) 的 数量 。
#
# 把字符串 s 视作一个 pair 列表：
# [(s[0],s[n−1]),(s[1],s[n−2]),(s[2],s[n−3]),⋯,(s[n−1],s[0])]
# 只要这个 pair 列表是另一个字符串 t 的 pair 列表的前缀，那么 s 就是 t 的前后缀。


from typing import Iterable, List, Tuple


class SimpleTrieNode:
    __slots__ = "children", "endCount"

    def __init__(self):
        self.children = dict()
        self.endCount = 0


class Solution:
    def countPrefixSuffixPairs(self, words: List[str]) -> int:
        def makePair(s: str) -> Iterable[Tuple[str, str]]:
            return zip(s, reversed(s))

        res = 0
        root = SimpleTrieNode()
        for w in words:
            cur = root
            for p in makePair(w):
                if p not in cur.children:
                    cur.children[p] = SimpleTrieNode()
                cur = cur.children[p]
                res += cur.endCount
            cur.endCount += 1
        return res
