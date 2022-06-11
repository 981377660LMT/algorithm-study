from collections import defaultdict
from typing import Dict, List


class TrieNode:
    __slots__ = ('wordCount', 'preCount', 'children')

    def __init__(self):
        self.wordCount = 0
        self.children: Dict[str, TrieNode] = defaultdict(TrieNode)


class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        if not word:
            return
        node = self.root
        for char in word:
            node = node.children[char]
        node.wordCount += 1

    def search(self, word: str) -> bool:
        if not word:
            return False
        node = self.root
        for char in word:
            if char not in node.children:
                return False
            node = node.children[char]
            if node.wordCount == 0:
                return False
        return True


# 找出 words 中所有的`前缀`都在 words 中的最长字符串。
# 1 <= words.length <= 105
class Solution:
    def longestWord(self, words: List[str]) -> str:
        trie = Trie()
        for word in words:
            trie.insert(word)

        words.sort(key=lambda x: (-len(x), x))
        for word in words:
            if trie.search(word):
                return word
        return ''


print(Solution().longestWord(words=["a", "banana", "app", "appl", "ap", "apply", "apple"]))
# 输出： "apple"
# 解释： "apple" 和 "apply" 都在 words 中含有各自的所有前缀。
# 然而，"apple" 在字典序中更小，所以我们返回之。

