# 3628. 插入一个字母的最大子序列数
# https://leetcode.cn/problems/maximum-number-of-subsequences-after-one-inserting/description/
# 你可以在字符串的 任意 位置（包括字符串的开头或结尾）最多插入一个 大写英文字母。
# 返回在 最多插入一个字母 后，字符串中可以形成的 "LCT" 子序列的 最大 数量。


class Solution:
    def numOfSubsequences(self, s: str) -> int:
        n = len(s)
        preL = [0] * (n + 1)
        for i, c in enumerate(s):
            preL[i + 1] = preL[i] + (c == "L")

        sufT = [0] * (n + 1)
        for i in range(n - 1, -1, -1):
            sufT[i] = sufT[i + 1] + (s[i] == "T")

        sufCT = [0] * (n + 1)
        for i in range(n - 1, -1, -1):
            sufCT[i] = sufCT[i + 1]
            if s[i] == "C":
                sufCT[i] += sufT[i + 1]

        preLC = [0] * (n + 1)
        countLCT = 0
        for i, c in enumerate(s):
            preLC[i + 1] = preLC[i]
            if c == "C":
                preLC[i + 1] += preL[i]
                countLCT += preL[i] * sufT[i + 1]

        best = 0
        for pos in range(n + 1):
            addL = sufCT[pos]
            addC = preL[pos] * sufT[pos]
            addT = preLC[pos]
            best = max(best, addL, addC, addT)
        return countLCT + best
