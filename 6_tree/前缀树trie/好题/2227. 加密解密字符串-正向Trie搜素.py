from collections import defaultdict
from typing import DefaultDict, List, Set


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

    def count(self, word: str, valueToKey: DefaultDict[str, Set[str]]) -> int:
        """统计并返回可以由 word 解密得到且出现在 dictionary 中的(前缀树中的) 字符串数目"""

        def dfs(node: TrieNode, cur: str) -> int:
            if not cur:
                return node.count

            prefix = cur[:2]
            chars = valueToKey[prefix]

            res = 0
            for char in chars:
                if char not in node.children:
                    continue
                res += dfs(node.children[char], cur[2:])
            return res

        return dfs(self.root, word)


class Encrypter:
    def __init__(self, keys: List[str], values: List[str], dictionary: List[str]):
        self.keyToValue = defaultdict(str)
        self.valueToKey = defaultdict(set)
        for k, v in zip(keys, values):
            self.keyToValue[k] = v
            self.valueToKey[v].add(k)

        self.trie = Trie()
        for word in dictionary:
            self.trie.insert(word)

    def encrypt(self, word1: str) -> str:
        return ''.join(self.keyToValue[c] for c in word1)

    def decrypt(self, word2: str) -> int:
        """
        统计并返回可以由 word2 解密得到且出现在 dictionary 中的字符串数目
        """
        return self.trie.count(word2, self.valueToKey)


# Your Encrypter object will be instantiated and called as such:
# obj = Encrypter(keys, values, dictionary)
# param_1 = obj.encrypt(word1)
# param_2 = obj.decrypt(word2)
