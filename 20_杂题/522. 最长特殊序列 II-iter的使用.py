from typing import List

# 最长特殊序列定义如下：该序列为某字符串独有的最长子序列（即不能是其他字符串的子序列）。

# 如果存在这样的特殊序列，那么它一定是某个给定的字符串。


class Solution:
    def findLUSlength(self, strs: List[str]) -> int:
        def isSubSequence(source: str, target: str) -> bool:
            it = iter(target)
            return all(sChar in it for sChar in source)

        for s in sorted(strs, key=len, reverse=True):
            print(s)
            if sum(isSubSequence(s, t) for t in strs) == 1:
                return len(s)
        return -1


print(Solution().findLUSlength(["aabbcc", "aabbcc", "cb", "abc"]))

# 3
