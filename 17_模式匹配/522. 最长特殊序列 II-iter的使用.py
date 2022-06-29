from typing import List

# 最长特殊序列定义如下：该序列为某字符串独有的最长子序列（即不能是其他字符串的子序列）。

# 如果存在这样的特殊序列，那么它一定是某个给定的字符串。


class Solution:
    def findLUSlength(self, strs: List[str]) -> int:
        def isSubSequence(longer: str, shorter: str) -> bool:
            if len(shorter) > len(longer):
                return False
            it = iter(longer)
            return all(need in it for need in shorter)

        for cur in sorted(strs, key=len, reverse=True):
            if sum(isSubSequence(other, cur) for other in strs) == 1:
                return len(cur)
        return -1


print(Solution().findLUSlength(["aabbcc", "aabbcc", "cb", "abc"]))

