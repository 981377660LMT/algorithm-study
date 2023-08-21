from collections import defaultdict
from bisect import bisect_right
from typing import List, Tuple


# 392. 判断子序列
def isSubSequence(longer: str, shorter: str) -> bool:
    if len(shorter) > len(longer):
        return False
    it = iter(longer)
    return all(need in it for need in shorter)


class Solution:
    def isSubsequence(self, shorter: str, longer: str) -> bool:
        """check if shorter is a subsequence of longer"""
        it = iter(shorter)
        return all(char in it for char in longer)

    # 如果有大量输入的 S，称作 S1, S2, ... , Sk 其中 k >= 10亿，你需要依次检查它们是否为 T 的子序列。在这种情况下，你会怎样改变代码？
    # 方法是预处理T使用邻接表存每种字符出现的索引
    def isSubsequence2(self, shorter: str, longer: str) -> bool:
        indexMap = defaultdict(list)
        for i, char in enumerate(longer):
            indexMap[char].append(i)

        pre = -1
        for char in shorter:
            if char not in indexMap:
                return False
            indexes = indexMap[char]
            nextPos = bisect_right(indexes, pre)
            if nextPos == len(indexes):
                return False
            pre = indexes[nextPos]

        return True

    def isSubsequence3(self, shorter: str, longer: str) -> bool:
        """子序列自动机"""

        n = len(longer)
        if n < len(shorter):
            return False
        nexts: List[Tuple[int, ...]] = [()] * n
        last = [-1] * 26
        for i in range(n - 1, -1, -1):
            last[ord(longer[i]) - 97] = i
            nexts[i] = tuple(last)

        nextIndex = 0
        for j in range(len(shorter)):
            if nextIndex >= n:
                return False
            cur = nexts[nextIndex]
            if cur[ord(shorter[j]) - 97] == -1:
                return False
            nextIndex = cur[ord(shorter[j]) - 97] + 1
        return True


print(Solution().isSubsequence2(shorter="abc", longer="ahbgdc"))
print(Solution().isSubsequence3(shorter="acb", longer="ahbgdc"))
