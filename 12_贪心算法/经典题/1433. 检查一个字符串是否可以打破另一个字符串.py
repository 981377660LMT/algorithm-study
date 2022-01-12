from collections import Counter
from string import ascii_lowercase


class Solution:
    def checkIfCanBreak(self, s1: str, s2: str) -> bool:
        a1 = [0] * 26
        a2 = [0] * 26
        for c in s1:
            a1[ord(c) - ord('a')] += 1
        for c in s2:
            a2[ord(c) - ord('a')] += 1

        flag1 = True
        flag2 = True
        acc1, acc2 = 0, 0

        for x in range(26):
            acc1 += a1[x]
            acc2 += a2[x]
            if acc1 < acc2:
                flag1 = False
            if acc2 < acc1:
                flag2 = False
        return flag1 or flag2


print(Solution().checkIfCanBreak(s1="abc", s2="xya"))
# 输出：true
# 解释："ayx" 是 s2="xya" 的一个排列，"abc" 是字符串 s1="abc" 的一个排列，且 "ayx" 可以打破 "abc" 。
