# 给你一个字符串 s 和一个字符串列表 wordDict 作为字典。
# 请你判断是否可以利用字典中出现的单词拼接出 s 。

# 注意：不要求字典中出现的单词全部都使用，并且字典中的单词可以重复使用。
from functools import lru_cache
from typing import List


class TrieNode:
    __slots__ = ('count', 'children')

    def __init__(self):
        self.count = 0
        self.children = {}


class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        node = self.root
        for char in word:
            if char not in node.children:
                node.children[char] = TrieNode()
            node = node.children[char]
        node.count += 1

    def search(self, word: str) -> int:
        """是否存在word，返回个数"""
        node = self.root
        for char in word:
            if char not in node.children:
                return 0
            node = node.children[char]
        return node.count


class Solution:
    def wordBreak(self, s: str, wordDict: List[str]) -> bool:
        @lru_cache(None)
        def dfs(index: int) -> bool:
            if index >= n:
                return index == n

            res = False
            for next in range(index + 1, n + 1):
                if trie.search(s[index:next]) > 0:
                    res = res or dfs(next)
            return res

        ok = set(wordDict)
        trie = Trie()
        for word in ok:
            trie.insert(word)

        n = len(s)
        res = dfs(0)
        dfs.cache_clear()
        return res
