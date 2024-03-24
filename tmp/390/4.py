from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e18)

# 给你两个字符串数组 wordsContainer 和 wordsQuery 。

# 对于每个 wordsQuery[i] ，你需要从 wordsContainer 中找到一个与 wordsQuery[i] 有 最长公共后缀 的字符串。如果 wordsContainer 中有两个或者更多字符串有最长公共后缀，那么答案为长度 最短 的。如果有超过两个字符串有 相同 最短长度，那么答案为它们在 wordsContainer 中出现 更早 的一个。


# 请你返回一个整数数组 ans ，其中 ans[i]是 wordsContainer中与 wordsQuery[i] 有 最长公共后缀 字符串的下标。


def min2(a: int, b: int) -> int:
    return a if a < b else b


class TrieNode:
    __slots__ = "children", "minIndex", "minLen"

    def __init__(self) -> None:
        self.children = dict()
        self.minIndex = INF
        self.minLen = INF


class Trie:
    __slots__ = "root"

    def __init__(self) -> None:
        self.root = TrieNode()

    def insert(self, word: str, index: int) -> None:
        cur = self.root
        for char in word:
            if char not in cur.children:
                cur.children[char] = TrieNode()
            cur = cur.children[char]
            if len(word) < cur.minLen:
                cur.minIndex = index
                cur.minLen = len(word)
            elif len(word) == cur.minLen:
                cur.minIndex = min2(cur.minIndex, index)

    def query(self, word: str) -> int:
        cur = self.root
        res = -1
        for c in word:
            if c not in cur.children:
                return res
            cur = cur.children[c]
            res = cur.minIndex
        return res


class Solution:
    def stringIndices(self, wordsContainer: List[str], wordsQuery: List[str]) -> List[int]:
        trie = Trie()
        for i, word in enumerate(wordsContainer):
            trie.insert(word[::-1], i)
        lens = [len(word) for word in wordsContainer]
        minLen = min(lens)
        minLenIndex = lens.index(minLen)
        res = [0] * len(wordsQuery)
        for i, word in enumerate(wordsQuery):
            cur = trie.query(word[::-1])
            res[i] = cur if cur != -1 else minLenIndex
        return res


# wordsContainer = ["abcd","bcd","xbcd"], wordsQuery = ["cd","bcd","xyz"]
print(Solution().stringIndices(["abcd", "bcd", "xbcd"], ["cd", "bcd", "xyz"]))
