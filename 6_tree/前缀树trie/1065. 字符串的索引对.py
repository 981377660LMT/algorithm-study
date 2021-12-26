from typing import List


class Solution:
    def indexPairs(self, text: str, words: List[str]) -> List[List[int]]:
        trie = {}
        for w in words:
            root = trie
            for char in w:
                if char not in root:
                    root[char] = {}
                root = root[char]
            root["is_word"] = True

        res = []
        n = len(text)
        for left in range(len(text)):
            root = trie
            right = left
            while right < n and text[right] in root:
                root = root[text[right]]
                right += 1
                if "is_word" in root:
                    res.append([left, right - 1])
        return res


print(Solution().indexPairs("ababa", ["aba", "ab"]))
