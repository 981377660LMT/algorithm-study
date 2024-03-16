# 745. 前缀和后缀搜索
# https://leetcode.cn/problems/prefix-and-suffix-search/description/
# 维护两个trie
# 单词查找最坏时间复杂度为O(n)


from typing import List


class TrieNode:
    __slots__ = "children", "maxIndex"

    def __init__(self) -> None:
        self.children = dict()
        self.maxIndex = -1


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
            cur.maxIndex = max(cur.maxIndex, index)

    def query(self, word: str) -> int:
        cur = self.root
        for c in word:
            if c not in cur.children:
                return -1
            cur = cur.children[c]
        return cur.maxIndex


class WordFilter:
    __slots__ = "_trie"

    def __init__(self, words: List[str]):
        """使用词典中的单词 words 初始化对象。"""
        self._trie = Trie()
        for i, word in enumerate(words):
            target = word + "#" + word
            for j in range(len(word) + 1):
                self._trie.insert(target[j:], i)

    def f(self, pref: str, suff: str) -> int:
        """
        返回词典中具有前缀 prefix 和后缀 suff 的单词的下标。
        如果存在不止一个满足要求的下标，返回其中 最大的下标 。如果不存在这样的单词，返回 -1 。
        """
        target = suff + "#" + pref
        return self._trie.query(target)


if __name__ == "__main__":
    obj = WordFilter(["apple"])
    print(obj.f("a", "e"))
    print(obj.f("b", ""))
