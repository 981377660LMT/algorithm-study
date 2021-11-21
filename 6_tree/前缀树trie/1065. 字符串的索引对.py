from typing import List


class Solution:
    def indexPairs(self, text: str, words: List[str]) -> List[List[int]]:
        tree = {}
        for w in words:
            root = tree
            for char in w:
                if char not in root:
                    root[char] = {}
                root = root[char]
            root["is_word"] = True

        res = []
        n = len(text)
        for start, _ in enumerate(text):
            root = tree
            cur = start
            while cur < n and text[cur] in root:
                root = root[text[cur]]
                cur += 1
                if "is_word" in root:
                    res.append([start, cur - 1])
        return res


print(Solution().indexPairs("ababa", ["aba", "ab"]))
