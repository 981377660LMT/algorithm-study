# 524. 通过删除字母匹配到字典里最长单词
# https://leetcode.cn/problems/longest-word-in-dictionary-through-deleting/description/
# 给你一个字符串 s 和一个字符串数组 dictionary ，找出并返回 dictionary 中最长的字符串，该字符串可以通过删除 s 中的某些字符得到。
# 如果答案不止一个，返回长度最长且字母序最小的字符串。如果答案不存在，则返回空字符串。


from typing import List
from SubsequenceAutomaton import SubsequenceAutomaton1


class Solution:
    def findLongestWord(self, s: str, dictionary: List[str]) -> str:
        res = ""
        S = SubsequenceAutomaton1(s)
        for w in dictionary:
            if S.includes(w):
                if len(w) > len(res) or (len(w) == len(res) and w < res):
                    res = w
        return res
