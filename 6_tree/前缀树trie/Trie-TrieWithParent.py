from collections import defaultdict
from typing import Dict, Iterable, List, Optional


class TrieNodeWithParent:
    __slots__ = ("wordCount", "preCount", "children", "parent")

    def __init__(self, parent: Optional["TrieNodeWithParent"] = None):
        self.wordCount = 0  # 以当前节点为结尾的单词个数
        self.preCount = 0  # 以当前节点为前缀的单词个数
        self.parent = parent
        self.children: Dict[str, TrieNodeWithParent] = defaultdict(lambda: TrieNodeWithParent(self))


class TrieWithParent:
    __slots__ = "root"

    def __init__(self, words: Optional[Iterable[str]] = None):
        self.root = TrieNodeWithParent()
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

    def countWordAsPrefix(self, s: str) -> List[int]:
        """对s的每个非空前缀pre,返回trie中有多少个单词是pre的前缀"""
        node = self._find(s)
        if node is None:
            return []

        diff = []
        while node.parent is not None:
            diff.append(node.wordCount)  # !重复的单词统计为多次而不是1次
            node = node.parent

        res, cur = [], 0
        for i in range(len(diff) - 1, -1, -1):
            cur += diff[i]
            res.append(cur)
        return res

    def remove(self, s: str) -> None:
        """从前缀树中移除`1个`s 需要保证s在前缀树中"""
        if not s:
            return
        node = self.root
        for char in s:
            if char not in node.children:
                raise ValueError(f"word {s} not in trie")
            if node.children[char].preCount == 1:
                del node.children[char]
                return
            node = node.children[char]
            node.preCount -= 1
        node.wordCount -= 1

    def _find(self, s: str) -> Optional[TrieNodeWithParent]:
        """返回s所在结点"""
        if not s:
            return None
        node = self.root
        for char in s:
            if char not in node.children:
                return None
            node = node.children[char]
        return node


if __name__ == "__main__":
    MOD = 998244353
    INV2 = 499122177
    n = int(input())
    names = [input() for _ in range(n)]

    trie = TrieWithParent(names)
    for name in names:
        count1 = trie.countWordStartsWith(name)[-1]  # name是多少个单词的前缀
        count2 = trie.countWordAsPrefix(name)[-1]  # 有多少个单词是name的前缀
        other = n - 1 - count1 - count2  # 互相没有前缀关系
        print((1 + count2 + other * INV2) % MOD)
