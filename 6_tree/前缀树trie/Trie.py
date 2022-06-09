class TrieNode:
    __slots__ = ('count', 'preCount', 'children')

    def __init__(self):
        self.count = 0
        self.preCount = 0
        self.children = dict()


class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        if not word:
            return
        node = self.root
        for char in word:
            if char not in node.children:
                node.children[char] = TrieNode()
            node = node.children[char]
            node.preCount += 1
        node.count += 1

    def countWord(self, word: str) -> int:
        """是否存在word,返回个数"""
        if not word:
            return 0
        node = self.root
        for char in word:
            if char not in node.children:
                return 0
            node = node.children[char]
        return node.count

    def countPrefix(self, prefix: str) -> int:
        """是否存在以prefix为前缀的单词,返回个数"""
        if not prefix:
            return 0
        node = self.root
        for char in prefix:
            if char not in node.children:
                return 0
            node = node.children[char]
        return node.preCount

    def discard(self, word: str) -> None:
        if not word:
            return
        node = self.root
        for char in word:
            if char not in node.children:
                return
            node = node.children[char]
            node.preCount -= 1
        node.count -= 1
