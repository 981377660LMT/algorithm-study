from collections import defaultdict
from typing import Dict, Iterable, Optional


class TrieNodeWithParent:
    __slots__ = ("wordCount", "preCount", "children", "parent")

    def __init__(self, parent: Optional["TrieNodeWithParent"] = None):
        self.wordCount = 0  # 以当前节点为结尾的单词个数
        self.preCount = 0  # 以当前节点为前缀的单词个数
        self.parent = parent
        self.children: Dict[str, TrieNodeWithParent] = defaultdict(lambda: TrieNodeWithParent(self))


class TrieWithParent:
    __slots__ = "root"

    def __init__(self, iterable: Optional[Iterable[str]] = None):
        self.root = TrieNodeWithParent()
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

    def countWordAsPrefix(self, word: str) -> int:
        """树中有多少个单词是word的前缀"""
        node = self._find(word)
        if node is None:
            return 0

        res = 0
        while node.parent is not None:
            res += node.wordCount  # 重复的单词统计为多次而不是1次
            node = node.parent
        return res

    def _find(self, word: str) -> Optional[TrieNodeWithParent]:
        """返回word所在结点"""
        if not word:
            return None
        node = self.root
        for char in word:
            if char not in node.children:
                return None
            node = node.children[char]
        return node
