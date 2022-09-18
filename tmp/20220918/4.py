"""教训 字符串长度小于1000切片近似常数 可以考虑暴力切片存储字符串"""

from collections import defaultdict
from typing import Dict, Iterable, List, Optional


class TrieNode:
    __slots__ = ("wordCount", "preCount", "children")

    def __init__(self):
        self.wordCount = 0
        self.preCount = 0
        self.children: Dict[str, TrieNode] = defaultdict(TrieNode)


class Trie:
    __slots__ = ("root",)

    def __init__(self, iterable: Optional[Iterable[str]] = None):
        self.root = TrieNode()
        for word in iterable or ():
            self.insert(word)

    def insert(self, word: str) -> None:
        if not word:
            return
        node = self.root
        for char in word:
            node = node.children[char]
            node.preCount += 1
        node.wordCount += 1

    def countWord(self, word: str) -> int:
        """树中有多少个单词word"""
        if not word:
            return 0
        node = self.root
        for char in word:
            if char not in node.children:
                return 0
            node = node.children[char]
        return node.wordCount

    def countWordStartsWith(self, prefix: str) -> int:
        """树中有多少个以prefix为前缀的单词"""
        if not prefix:
            return 0
        node = self.root
        for char in prefix:
            if char not in node.children:
                return 0
            node = node.children[char]
        return node.preCount

    def remove(self, word: str) -> None:
        """从前缀树中移除`1个`word 需要保证word在前缀树中"""
        if not word:
            return
        node = self.root
        for char in word:
            if char not in node.children:
                raise ValueError(f"word {word} not in trie")
            node = node.children[char]
            node.preCount -= 1
        node.wordCount -= 1


class Solution:
    def sumPrefixScores(self, words: List[str]) -> List[int]:
        def dfs(root: "TrieNode", pre: str) -> None:
            for key, child in root.children.items():
                nextPre = pre + key
                mp[nextPre] = child.preCount
                dfs(child, nextPre)

        trie = Trie(words)
        mp = defaultdict(int)
        root = trie.root
        dfs(root, "")

        res = []
        for word in words:
            cur = 0
            for i in range(len(word)):
                pre = word[: i + 1]
                cur += mp[pre]
            res.append(cur)
        return res


print(Solution().sumPrefixScores(words=["abc", "ab", "bc", "b"]))
