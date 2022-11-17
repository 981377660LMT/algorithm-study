from bisect import bisect_right
from collections import defaultdict
from 子序列自动机 import SubsequenceAutomaton1


# 最小窗口子序列
class Solution:
    def minWindow2(self, s1: str, s2: str) -> str:
        SA = SubsequenceAutomaton1(s1)
        starts = [i for i, char in enumerate(s1) if char == s2[0]]
        res = None
        for start in starts:
            hit, end = SA.match(s2, start)
            if hit != len(s2):
                continue

            len_ = end - start + 1
            if res is None or len_ < res[1] - res[0] + 1:
                res = [start, end]

        return s1[res[0] : res[1] + 1] if res is not None else ""

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


print(Solution().minWindow2("abcdebdde", "bde"))
