# 2781. 最长合法子字符串的长度
# https://leetcode.cn/problems/length-of-the-longest-valid-substring/
# 给你一个字符串 word 和一个字符串数组 forbidden 。
# 如果一个字符串不包含 forbidden 中的任何字符串，我们称这个字符串是 合法 的。
# 请你返回字符串 word 的一个 最长合法子字符串 的长度。
# 子字符串 指的是一个字符串中一段连续的字符，它可以为空。
#
# 1 <= word.length <= 1e5
# word 只包含小写英文字母。
# 1 <= forbidden.length <= 1e5
# !1 <= forbidden[i].length <= 1e5
# !sum(len(forbidden)) <= 1e7
# forbidden[i] 只包含小写英文字母。
#
# 思路:
# 类似字符流, 需要处理出每个位置为结束字符的包含至少一个模式串的`最短后缀`.
# !那么此时左端点就对应这个位置+1


from collections import defaultdict
from typing import List
from ACAutoMatonLegacy import ACAutoMatonLegacy

INF = int(1e18)


class Solution:
    def longestValidSubstring(self, word: str, forbidden: List[str]) -> int:
        def didInsert(pos: int) -> None:
            minLen[pos] = min(minLen[pos], len(s))

        def dp(pre: int, cur: int) -> None:
            minLen[cur] = min(minLen[cur], minLen[pre])

        minLen = defaultdict(lambda: INF)  # !ac自动机每个状态匹配到的模式串的最小长度
        acm = ACAutoMatonLegacy()
        for i, s in enumerate(forbidden):
            acm.insert(i, s, didInsert=didInsert)
        acm.build(dp=dp)

        res = 0
        left = 0
        state = 0
        for right, char in enumerate(word):
            state = acm.move(state, char)
            left = max(left, right - minLen[state] + 2)
            res = max(res, right - left + 1)
        return res
