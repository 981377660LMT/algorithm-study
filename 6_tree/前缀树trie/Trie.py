class TrieNode:
    __slots__ = ('count', 'preCount', 'children')

    def __init__(self):
        self.count = 0
        self.preCount = 0
        self.children = {}


class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        node = self.root
        for char in word:
            if char not in node.children:
                node.children[char] = TrieNode()
            node = node.children[char]
            node.preCount += 1
        node.count += 1

    def search(self, word: str) -> bool:
        node = self.root
        for char in word:
            if char not in node.children:
                return False
            node = node.children[char]
        return node.count > 0

    def startsWith(self, prefix: str) -> int:
        node = self.root
        for char in prefix:
            if char not in node.children:
                return False
            node = node.children[char]
        return node.preCount > 0
