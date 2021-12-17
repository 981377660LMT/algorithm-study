from collections import Counter
from itertools import groupby

# 1 <= text.length <= 20000
# 给你一个字符串 text，你只能交换其中两个字符一次或者什么都不做，然后得到一些单字符重复的子串。返回其中最长的子串的长度。


# [k for k, g in groupby('AAAABBBCCDAABBB')] --> A B C D A B
# [list(g) for k, g in groupby('AAAABBBCCD')] --> AAAA BBB CC D


# 考虑两种情况
# 1. ..aaaabaaaa.. 被一个b分隔了，把b换成a   `remove 1 divider`
# 2. ..aaaa...aa.. => 从后面拿一个a过来    `extend 1`
class Solution:
    def maxRepOpt1(self, text: str) -> int:
        counter = Counter(text)
        groups = [[char, len(list(group))] for char, group in groupby(text)]
        print(groups)

        # 1. extend 1 情形
        res = max(min(count + 1, counter[char]) for char, count in groups)

        # 2. remove 1 divider 情形
        for i in range(1, len(groups) - 1):
            if groups[i - 1][0] == groups[i + 1][0] and groups[i][1] == 1:
                sameChar = groups[i - 1][0]
                res = max(res, min(groups[i - 1][1] + groups[i + 1][1] + 1, counter[sameChar]))

        return res


print(Solution().maxRepOpt1(text="aaabaaa"))
# 输出：6
