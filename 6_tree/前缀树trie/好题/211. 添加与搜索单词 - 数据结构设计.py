from collections import defaultdict


class TrieNode:
    __slots__ = ('count', 'preCount', 'children')

    def __init__(self):
        self.count = 0
        self.preCount = 0
        self.children = defaultdict(TrieNode)


class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        node = self.root
        for char in word:
            node = node.children[char]
            node.preCount += 1
        node.count += 1

    def startsWith(self, prefix: str) -> int:
        """是否存在以prefix为前缀的单词，返回个数"""
        node = self.root
        for char in prefix:
            if char not in node.children:
                return 0
            node = node.children[char]
        return node.preCount

    def delete(self, word: str) -> None:
        node = self.root
        for char in word:
            if char not in node.children:
                return
            node = node.children[char]
            node.preCount -= 1
        node.count -= 1

    def search(self, word: str) -> bool:
        """字典中是否存在word
        word 中可能包含一些 '.' ，每个 . 都可以表示任何一个字母
        """

        def dfs(node: TrieNode, cur: str) -> bool:
            # 可以用index优化，而非切片
            if not cur:
                return node.count > 0
            if cur[0] == '.':
                for child in node.children.values():
                    if dfs(child, cur[1:]):
                        return True
                return False
            else:
                if cur[0] not in node.children:
                    return False
                return dfs(node.children[cur[0]], cur[1:])

        return dfs(self.root, word)


class WordDictionary:
    """请你设计一个数据结构，支持 添加新单词 和 查找字符串是否与任何先前添加的字符串匹配
    """

    def __init__(self):
        self.trie = Trie()

    def addWord(self, word: str) -> None:
        """将 word 添加到数据结构中，之后可以对它进行匹配"""
        self.trie.insert(word)

    def search(self, word: str) -> bool:
        return self.trie.search(word)
