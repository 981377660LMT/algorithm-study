from typing import List
from collections import Counter


# counter配对题
# 这题卡了很久，是因为还没想清楚就开始写了，然后debug好久；以后一定要先把情况考虑清楚，丢到counter里辅助思考

# 情况就是配对+一个same
class Solution:
    def longestPalindrome(self, words: List[str]) -> int:
        same = Counter()
        diff = Counter()

        for word in words:
            if word[0] == word[1]:
                same[word] += 1
            else:
                diff[word] += 1

        res = 0
        for word in diff:
            match = word[::-1]
            if match in diff:
                res += min(diff[word], diff[match]) * 2
        for word in same:
            res += same[word] // 2 * 4
            same[word] %= 2

        for word in same:
            if same[word] == 1:
                res += 2
                break

        return res


# 6 8 2 14 22 14
print(Solution().longestPalindrome(words=["lc", "cl", "gg"]))
print(Solution().longestPalindrome(words=["ab", "ty", "yt", "lc", "cl", "ab"]))
print(Solution().longestPalindrome(words=["cc", "ll", "xx"]))
print(
    Solution().longestPalindrome(
        words=[
            "qo",
            "fo",
            "fq",
            "qf",
            "fo",
            "ff",
            "qq",
            "qf",
            "of",
            "of",
            "oo",
            "of",
            "of",
            "qf",
            "qf",
            "of",
        ]
    )
)
print(
    Solution().longestPalindrome(
        ["dd", "aa", "bb", "dd", "aa", "dd", "bb", "dd", "aa", "cc", "bb", "cc", "dd", "cc"]
    )
)
print(
    Solution().longestPalindrome(words=["em", "pe", "mp", "ee", "pp", "me", "ep", "em", "em", "me"])
)
