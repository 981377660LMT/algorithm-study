from typing import List


class Solution:
    def longestWord(self, words: List[str]) -> str:
        """
        给出一个字符串数组 words 组成的一本英语词典。返回 words 中最长的一个单词，
        该单词是由 words 词典中其他单词逐步添加一个字母组成。
        """
        words = sorted(words)
        dp = set()
        res = ""
        for w in words:
            # w[:-1] in cur 此处可用前缀树优化
            if len(w) == 1 or w[:-1] in dp:
                dp.add(w)
                if len(w) > len(res):
                    res = w
        return res


print(Solution().longestWord(words=["a", "banana", "app", "appl", "ap", "apply", "apple"]))
