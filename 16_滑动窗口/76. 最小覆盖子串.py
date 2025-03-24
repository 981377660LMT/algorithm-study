# 76. 最小覆盖子串
# https://leetcode.cn/problems/minimum-window-substring/description/?envType=problem-list-v2&envId=e2xDEg9M


# !给你一个字符串 s 、一个字符串 t 。返回 s 中涵盖 t 所有字符的最小子串。如果 s 中不存在涵盖 t 所有字符的子串，则返回空字符串 "" 。
class Solution:
    def minWindow(self, s: str, t: str) -> str:
        if len(s) < len(t):
            return ""
