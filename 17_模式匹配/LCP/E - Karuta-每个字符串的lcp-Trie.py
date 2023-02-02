"""
给定n个字符串,对于每个字符串 si,问 maxLCP(si,sj)i≠j
其中LCP是最长公共前缀。

trie 求lcp
n<=5e5 所有字符串总长<=5e5
"""


from typing import List


def karuta(words: List[str]) -> List[int]:
    n = len(words)
    trie = Trie(words)
    res = [0] * n
    for i, w in enumerate(words):
        trie.remove(w)
        res[i] = trie.lcp(w)
        trie.insert(w)
    return res


if __name__ == "__main__":
    n = int(input())
    words = [input() for _ in range(n)]

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

        def lcp(self, s: str) -> int:
            """返回s在trie中的最长公共前缀"""
            if not s:
                return 0
            node = self.root
            for i, char in enumerate(s):
                if char not in node.children:
                    return i
                node = node.children[char]
            return len(s)

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
                if node.children[char].preCount == 1:
                    del node.children[char]
                    return
                node = node.children[char]
                node.preCount -= 1
            node.wordCount -= 1

    print(*karuta(words), sep="\n")
