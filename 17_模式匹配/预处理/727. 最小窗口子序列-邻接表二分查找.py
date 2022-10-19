from bisect import bisect_right
from collections import defaultdict

# 最小窗口子序列
class Solution:
    def minWindow(self, s1: str, s2: str) -> str:
        """
        找出 s1 中最短的（连续）子串 W ，使得 s2 是 W 的 子序列 。
        如果有不止一个最短长度的窗口，返回开始位置最靠左的那个。
        """

        indexMap = defaultdict(list)
        s2Set = set(s2)
        for i, char in enumerate(s1):
            if char in s2Set:
                indexMap[char].append(i)

        res = ""
        # 枚举出发点
        for start in indexMap[s2[0]]:
            cur = start
            for char in s2[1:]:
                index = bisect_right(indexMap[char], cur)
                if index == len(indexMap[char]):
                    break
                cur = indexMap[char][index]
            else:
                if not res or cur - start + 1 < len(res):
                    res = s1[start : cur + 1]

        return res
