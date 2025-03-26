# 76. 最小覆盖子串
# https://leetcode.cn/problems/minimum-window-substring/description/?envType=problem-list-v2&envId=e2xDEg9M


# !给你一个字符串 s 、一个字符串 t 。返回 s 中涵盖 t 所有字符的最小子串。
# 如果 s 中不存在涵盖 t 所有字符的子串，则返回空字符串 "" 。


from collections import defaultdict


class Solution:
    def minWindow(self, s: str, t: str) -> str:
        if len(s) < len(t):
            return ""

        resLeft, resLen = -1, len(s) + 1
        counter = defaultdict(int)
        for c in t:
            counter[c] += 1
        less = len(counter)  # 未覆盖的字符种类数

        left = 0
        for right, c in enumerate(s):
            counter[c] -= 1
            if counter[c] == 0:
                less -= 1
            while left <= right and less == 0:
                if right - left + 1 < resLen:
                    resLeft, resLen = left, right - left + 1
                x = s[left]
                if counter[x] == 0:
                    less += 1
                counter[x] += 1
                left += 1

        return "" if resLeft == -1 else s[resLeft : resLeft + resLen]
