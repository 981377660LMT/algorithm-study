from typing import List, Tuple


MOD = int(1e9 + 7)


class Solution:
    def maximumSubsequenceCount(self, text: str, pattern: str) -> int:
        def getCount(t: str):
            if pattern[0] == pattern[1]:
                count = t.count(pattern[0])
                return count * (count - 1) // 2
            else:
                p1, p2 = 0, 0
                for char in t:
                    if char == pattern[0]:
                        p1 += 1
                    elif char == pattern[1]:
                        p2 += p1
                return p2

        text1, text2 = pattern[0] + text, text + pattern[1]

        return max(getCount(text1), getCount(text2))


print(Solution().maximumSubsequenceCount("fwyymvreuftzgrcrxczjacqovduqaiig", "yy"))
