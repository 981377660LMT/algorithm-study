from collections import defaultdict
from typing import Dict, Iterable, List, Optional


class TrieNode:
    __slots__ = ("wordCount", "preCount", "children")

    def __init__(self):
        self.wordCount = 0
        self.preCount = 0
        self.children: Dict[str, TrieNode] = defaultdict(TrieNode)


class Trie:
    __slots__ = "root"

    def __init__(self, words: Optional[Iterable[str]] = None):
        self.root = TrieNode()
        for word in words or ():
            self.insert(word)

    def insert(self, s: str) -> None:
        if not s:
            return
        node = self.root
        for char in s:
            node = node.children[char]
            node.preCount += 1
        node.wordCount += 1

    def countWord(self, s: str) -> List[int]:
        """对s的每个非空前缀pre,返回trie中有多少个等于pre的单词"""
        if not s:
            return []
        res = []
        node = self.root
        for char in s:
            if char not in node.children:
                return []
            node = node.children[char]
            res.append(node.wordCount)
        return res

    def countWordStartsWith(self, s: str) -> List[int]:
        """对s的每个非空前缀pre,返回trie中有多少个单词以pre为前缀"""
        if not s:
            return []
        res = []
        node = self.root
        for char in s:
            if char not in node.children:
                return []
            node = node.children[char]
            res.append(node.preCount)
        return res

    def remove(self, s: str) -> None:
        """从前缀树中移除`1个`s 需要保证s在前缀树中"""
        if not s:
            return
        node = self.root
        for char in s:
            if char not in node.children:
                raise ValueError(f"word {s} not in trie")
            node = node.children[char]
            node.preCount -= 1
        node.wordCount -= 1
