# 745. 前缀和后缀搜索
# https://leetcode.cn/problems/prefix-and-suffix-search/description/
# 维护两个trie
# 单词查找最坏时间复杂度为O(n)


from typing import List


class TrieNode:
    __slots__ = "children", "indexes"

    def __init__(self) -> None:
        self.children = dict()
        self.indexes = []  # 记录下标


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
            cur.indexes.append(index)

    def query(self, word: str) -> List[int]:
        cur = self.root
        for c in word:
            if c not in cur.children:
                return []
            cur = cur.children[c]
        return cur.indexes


class WordFilter:
    __slots__ = "_prefixTrie", "_suffixTrie"

    def __init__(self, words: List[str]):
        """使用词典中的单词 words 初始化对象。"""
        self._prefixTrie = Trie()
        self._suffixTrie = Trie()
        for i, word in enumerate(words):
            self._prefixTrie.insert(word, i)
            self._suffixTrie.insert(word[::-1], i)

    def f(self, pref: str, suff: str) -> int:
        """
        返回词典中具有前缀 prefix 和后缀 suff 的单词的下标。
        如果存在不止一个满足要求的下标，返回其中 最大的下标 。如果不存在这样的单词，返回 -1 。
        """
        arr1 = self._prefixTrie.query(pref)
        if not arr1:
            return -1
        arr2 = self._suffixTrie.query(suff[::-1])
        if not arr2:
            return -1
        i, j = len(arr1) - 1, len(arr2) - 1  # 找最大的下标
        while i >= 0 and j >= 0:
            if arr1[i] == arr2[j]:
                return arr1[i]
            elif arr1[i] > arr2[j]:
                i -= 1
            else:
                j -= 1
        return -1


if __name__ == "__main__":
    obj = WordFilter(["apple"])
    print(obj.f("a", "e"))
    print(obj.f("b", ""))
