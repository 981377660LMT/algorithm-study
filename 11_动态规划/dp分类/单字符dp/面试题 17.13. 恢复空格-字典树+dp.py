from functools import lru_cache
from typing import List

# 把文章断开，要求未识别的字符最少，返回未识别的字符数。
# 0 <= len(sentence) <= 1000
# dictionary中总字符数不超过 150000。

# 使用字典树可利用字符串的公共前缀来减少查询时间，最大限度地减少无谓的字符串比较


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
    def respace(self, dictionary: List[str], sentence: str) -> int:
        # 时间复杂度O(|dictionary|+O(n^2))
        # 建字典树的时间复杂度取决于单词的总字符数
        # dp转移时间复杂度O(n^2)
        @lru_cache(None)
        def dfs(index: int) -> int:
            if index >= n:
                return 0 if index == n else int(1e20)

            res = 1 + dfs(index + 1)
            for next in range(index + 1, n + 1):  # !根据index转移
                if trie.search(sentence[index:next]) > 0:
                    res = min(res, dfs(next))
            return res

            for cur in ok:  # !根据单词转移
                len_ = len(cur)
                if cur == sentence[index : index + len_]:  # 可用字符串哈希优化
                    res = min(res, dfs(index + len_))
            return res

        ok = set(dictionary)
        trie = Trie()
        for word in ok:
            trie.insert(word)

        n = len(sentence)
        res = dfs(0)
        dfs.cache_clear()
        return res
