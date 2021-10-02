from typing import List
from collections import Counter


class Solution:
    def shortestCompletingWord(self, licensePlate: str, words: List[str]) -> str:
        need = Counter(c for c in licensePlate.lower() if c.isalpha())
        words.sort(key=len)  # 从短的找起
        for word in words:
            asset = Counter(c for c in word.lower() if c.isalpha())
            print(need - asset)
            if not (need - asset):  # asset完全包含了need, 输出空的Counter()
                return word
        return ''


print(Solution().shortestCompletingWord("1s3 PSt", ["step", "steps", "stripe", "stepple"]))
# 题目数据保证一定存在一个最短补全词。当有多个单词都符合最短补全词的匹配条件时取单词列表中最靠前的一个。
