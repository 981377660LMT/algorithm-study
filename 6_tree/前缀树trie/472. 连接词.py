from typing import List
from collections import defaultdict
from typing import Dict


class TrieNode:
    __slots__ = ('wordCount', 'preCount', 'children')

    def __init__(self):
        self.wordCount = 0
        self.children: Dict[str, TrieNode] = defaultdict(TrieNode)


class Trie:
    def __init__(self):
        self.root = TrieNode()  # 声明为公有 便于改造

    def insert(self, word: str) -> None:
        if not word:
            return
        node = self.root
        for char in word:
            node = node.children[char]
        node.wordCount += 1

    def search(self, word: str) -> bool:
        def dfs(root: TrieNode, word: str) -> bool:
            """word是由给定数组中的至少两个较短单词组成的字符串"""
            node = root
            for i in range(len(word)):
                if word[i] not in node.children:
                    return False
                node = node.children[word[i]]
                if node.wordCount > 0:
                    if dfs(self.root, word[i + 1 :]):
                        return True
            return node.wordCount > 0

        return dfs(self.root, word)


# 1 <= words.length <= 1e4
# 0 <= words[i].length <= 30
# 0 <= sum(words[i].length) <= 1e5
class Solution:
    def findAllConcatenatedWordsInADict(self, words: List[str]) -> List[str]:
        trie, res = Trie(), []
        for w in sorted(words, key=len):
            if not w:
                continue
            if trie.search(w):
                res.append(w)
            trie.insert(w)
        return res


print(
    Solution().findAllConcatenatedWordsInADict(
        ["cat", "cats", "catsdogcats", "dog", "dogcatsdog", "hippopotamuses", "rat", "ratcatdogcat"]
    )
)

# 输出：["catsdogcats","dogcatsdog","ratcatdogcat"]
# 解释："catsdogcats" 由 "cats", "dog" 和 "cats" 组成;
#      "dogcatsdog" 由 "dog", "cats" 和 "dog" 组成;
#      "ratcatdogcat" 由 "rat", "cat", "dog" 和 "cat" 组成。
