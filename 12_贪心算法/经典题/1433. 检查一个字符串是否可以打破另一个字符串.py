from collections import Counter
from string import ascii_lowercase


class Solution:
    def checkIfCanBreak(self, s1: str, s2: str) -> bool:
        def check(c1, c2):
            diff = 0
            for char in ascii_lowercase:
                diff += c1[char] - c2[char]
                if diff < 0:
                    return False
            return True

        c1, c2 = Counter(s1), Counter(s2)
        return check(c1, c2) or check(c2, c1)


print(Solution().checkIfCanBreak(s1="abc", s2="xya"))
# 输出：true
# 解释："ayx" 是 s2="xya" 的一个排列，"abc" 是字符串 s1="abc" 的一个排列，且 "ayx" 可以打破 "abc" 。
