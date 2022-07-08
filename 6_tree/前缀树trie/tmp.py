from typing import Iterable, List, Optional
from collections import defaultdict
from typing import Dict


class Solution:
    def replaceWords(self, dictionary: List[str], sentence: str) -> str:
        """你需要将句子中的所有继承词用词根替换掉。如果继承词有许多可以形成它的词根，则用最短的词根替换它。"""
        words, trie = sentence.split(), Trie(dictionary)
        return " ".join(trie.search(w) for w in words)


class TrieNode:
    __slots__ = ("wordCount", "preCount", "children", "word")

    def __init__(self):
        self.wordCount = 0
        self.preCount = 0
        self.word = ""  # 存储当前的单词
        self.children: Dict[str, TrieNode] = defaultdict(TrieNode)


class Trie:
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
        node.word = word

    def search(self, word: str) -> str:
        if not word:
            return word
        node = self.root
        for char in word:
            if char not in node.children:
                return word
            node = node.children[char]
            if node.wordCount > 0:
                return node.word
        return word


print(
    Solution().replaceWords(
        ["cat", "bat", "rat"], "the cattle was rattled by the battery"
    )
)
