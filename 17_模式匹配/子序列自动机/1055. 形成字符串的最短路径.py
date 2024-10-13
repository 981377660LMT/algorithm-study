# 1055. 形成字符串的最短路径
# https://leetcode.cn/problems/shortest-way-to-form-string/description/
#
# 对于任何字符串，我们可以通过删除其中一些字符（也可能不删除）来构造该字符串的 子序列 。
# (例如，“ace” 是 “abcde” 的子序列，而 “aec” 不是)。
#
# 给定源字符串 source 和目标字符串 target，
# 返回 源字符串 source 中能通过串联形成目标字符串 target 的 子序列 的最小数量 。
# 如果无法通过串联源字符串中的子序列来构造目标字符串，则返回 -1。

from SubsequenceAutomaton import SubsequenceAutomaton1


class Solution:
    def shortestWay(self, source: str, target: str) -> int:
        S = SubsequenceAutomaton1(source)
        res = 0
        ti = 0
        while True:
            hit, _ = S.match(target, tStart=ti)
            if hit == 0:
                return -1
            res += 1
            ti += hit
            if ti == len(target):
                return res
