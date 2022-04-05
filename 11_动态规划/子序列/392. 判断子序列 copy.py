from collections import defaultdict
from bisect import bisect_right


class Solution:
    def isSubsequence(self, s: str, t: str) -> bool:
        """check if s is a subsequence of t"""
        it = iter(t)
        return all(char in it for char in s)

    # 如果有大量输入的 S，称作 S1, S2, ... , Sk 其中 k >= 10亿，你需要依次检查它们是否为 T 的子序列。在这种情况下，你会怎样改变代码？

    # 方法是预处理T使用邻接表存每种字符出现的索引
    def isSubsequence2(self, s: str, t: str) -> bool:
        indexes = defaultdict(list)
        for i, char in enumerate(t):
            indexes[char].append(i)

        pre = -1
        for char in s:
            if char not in indexes:
                return False
            cand = indexes[char]
            nextPos = bisect_right(cand, pre)
            if nextPos == len(cand):
                return False
            pre = cand[nextPos]

        return True


print(Solution().isSubsequence2(s="abc", t="ahbgdc"))
