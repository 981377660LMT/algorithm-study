from typing import Generator, Iterable, Optional, Tuple


class TrieNode:
    __slots__ = ("wordCount", "preCount", "children")

    def __init__(self):
        self.wordCount = 0
        self.preCount = 0
        self.children = dict()


class Trie:
    __slots__ = "root"

    def __init__(self, words: Optional[Iterable[str]] = None):
        self.root = TrieNode()
        for word in words or ():
            self.insert(word)

    def insert(self, s: str) -> "TrieNode":
        if not s:
            return self.root
        node = self.root
        for char in s:
            if char not in node.children:
                newNode = TrieNode()
                node.children[char] = newNode
                node = newNode
            else:
                node = node.children[char]
            node.preCount += 1
        node.wordCount += 1
        return node

    def remove(self, s: str) -> "TrieNode":
        """从前缀树中移除`1个`s 需要保证s在前缀树中"""
        if not s:
            return self.root
        node = self.root
        for char in s:
            node = node.children[char]
            node.preCount -= 1
        node.wordCount -= 1
        return node

    def find(self, s: str) -> Optional[TrieNode]:
        """返回s所在结点"""
        if not s:
            return None
        node = self.root
        for char in s:
            if char not in node.children:
                return None
            node = node.children[char]
        return node

    def enumerate(self, s: str) -> Generator[Tuple[int, TrieNode], None, None]:
        if not s:
            return
        node = self.root
        for i, char in enumerate(s):
            if char not in node.children:
                return
            node = node.children[char]
            yield i, node


if __name__ == "__main__":

    class Trie2:
        def __init__(self):
            self.trie = Trie()

        def insert(self, word: str) -> None:
            self.trie.insert(word)

        def countWordsEqualTo(self, word: str) -> int:
            for i, node in self.trie.enumerate(word):
                if i == len(word) - 1:
                    return node.wordCount
            return 0

        def countWordsStartingWith(self, prefix: str) -> int:
            for i, node in self.trie.enumerate(prefix):
                if i == len(prefix) - 1:
                    return node.preCount
            return 0

        def erase(self, word: str) -> None:
            self.trie.remove(word)
