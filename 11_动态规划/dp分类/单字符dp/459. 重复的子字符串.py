# 459. 重复的子字符串
# 给定一个非空的字符串 s ，检查是否可以通过由它的一个子串重复多次(>=2)构成。
# https://leetcode.cn/problems/repeated-substring-pattern/description/


def minimalPeriod(s: str) -> int:
    """
    计算字符串 s 的最小周期。
    通过 (s + s).find(s, 1) 找到第一个非零的周期长度。
    如果没有找到，返回 len(s)。
    """
    n = len(s)
    res = (s + s).find(s, 1, -1)
    return res if res != -1 else n


class Solution:
    def repeatedSubstringPattern(self, s: str) -> bool:
        return minimalPeriod(s) != len(s)
