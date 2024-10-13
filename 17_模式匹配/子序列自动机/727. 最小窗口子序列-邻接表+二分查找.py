# https://leetcode.cn/problems/minimum-window-subsequence/
# 727. 最小窗口子序列
# 给定字符串 S and T，找出 S 中最短的（连续）子串 W ，使得 T 是 W 的 子序列 。
# !如果 S 中没有窗口可以包含 T 中的所有字符，返回空字符串 ""。如果有不止一个最短长度的窗口，返回开始位置最靠左的那个。

from SubsequenceAutomaton import SubsequenceAutomaton1, SubsequenceAutomaton2


class Solution:
    def minWindow(self, s1: str, s2: str) -> str:
        SA = SubsequenceAutomaton1(s1)
        starts = [i for i, v in enumerate(s1) if v == s2[0]]
        res = None
        for sStart in starts:
            hit, sEnd = SA.match(s2, sStart=sStart)
            if hit != len(s2):
                continue

            sLen = sEnd - sStart
            if res is None or sLen < res[1] - res[0]:
                res = [sStart, sEnd]

        return s1[res[0] : res[1]] if res is not None else ""


print(Solution().minWindow("abcdebdde", "bde"))
