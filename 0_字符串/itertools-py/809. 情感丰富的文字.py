from itertools import groupby
from typing import List


class Solution:
    def expressiveWords(self, s: str, words: List[str]) -> int:
        def check(word: str) -> bool:
            wordGroups = [[char, len(list(group))] for char, group in groupby(word)]
            return len(wordGroups) == len(groups) and all(
                char1 == char2 and (len1 == len2) or len2 <= len1 >= 3
                for (char1, len1), (char2, len2) in zip(groups, wordGroups)
            )

        groups = [[char, len(list(group))] for char, group in groupby(s)]
        return sum(check(w) for w in words)


print(Solution().expressiveWords(s="heeellooo", words=["hello", "hi", "helo"]))
# 输出：1
# 解释：
# 我们能通过扩张 "hello" 的 "e" 和 "o" 来得到 "heeellooo"。
# 我们不能通过扩张 "helo" 来得到 "heeellooo" 因为 "ll" 的长度小于 3 。
