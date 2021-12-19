from typing import List

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


class Trie:
    def __init__(self) -> None:
        self.root = {}

    def insert(self, word):
        root = self.root
        for c in word:
            if c not in root:
                root[c] = {}
            root = root[c]
        root['end'] = True

    def search(self, word):
        root = self.root
        for c in word:
            if c not in root:
                return False
            root = root[c]
            # 关键,要包含所有前缀
            if 'end' not in root:
                return False
        return True


print(Solution().longestWord(words=["a", "banana", "app", "appl", "ap", "apply", "apple"]))
# 输出： "apple"
# 解释： "apple" 和 "apply" 都在 words 中含有各自的所有前缀。
# 然而，"apple" 在字典序中更小，所以我们返回之。

