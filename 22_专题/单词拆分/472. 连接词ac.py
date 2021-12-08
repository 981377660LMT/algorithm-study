from typing import List
from functools import lru_cache


class Trie:
    def __init__(self):
        self.root = {}

    def insert(self, word):
        root = self.root
        for char in word:
            if char not in root:
                root[char] = {}
            root = root[char]
        root['#'] = 1

    @lru_cache(None)
    def countWords(self, word):
        if not word:
            return 0

        root = self.root
        res = float('-inf')

        for i, char in enumerate(word):
            if char not in root:
                return res

            root = root[char]
            if '#' in root:
                res = max(res, 1 + self.countWords(word[i + 1 :]))

        return res


# 1 <= words.length <= 104
# 0 <= words[i].length <= 1000
class Solution:
    def findAllConcatenatedWordsInADict(self, words: List[str]) -> List[str]:
        self.trie = Trie()
        res = []

        for word in words:
            self.trie.insert(word)
        for word in words:
            if self.trie.countWords(word) >= 2:
                res.append(word)
        return res


print(
    Solution().findAllConcatenatedWordsInADict(
        words=[
            "cat",
            "cats",
            "catsdogcats",
            "dog",
            "dogcatsdog",
            "hippopotamuses",
            "rat",
            "ratcatdogcat",
        ]
    )
)

# 输出：["catsdogcats","dogcatsdog","ratcatdogcat"]
# 解释："catsdogcats" 由 "cats", "dog" 和 "cats" 组成;
#      "dogcatsdog" 由 "dog", "cats" 和 "dog" 组成;
#      "ratcatdogcat" 由 "rat", "cat", "dog" 和 "cat" 组成。

