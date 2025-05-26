# 3557. 不相交子字符串的最大数量-划分型dp
# https://leetcode.cn/problems/find-maximum-number-of-non-intersecting-substrings/solutions/3685351/hua-fen-xing-tan-xin-fu-ti-dan-pythonjav-9i8j/
# 给你一个字符串 word。
# 返回以 首尾字母相同 且 长度至少为 4 的 不相交子字符串 的最大数量。
# 子字符串 是字符串中连续的 非空 字符序列。


class Solution:
    def maxSubstrings(self, word: str) -> int:
        res = 0
        firstPos = dict()
        for i, c in enumerate(word):
            if c not in firstPos:
                firstPos[c] = i
            elif i - firstPos[c] >= 3:
                res += 1
                firstPos.clear()  # 找下一个子串
        return res
