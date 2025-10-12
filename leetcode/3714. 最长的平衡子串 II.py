# 3714. 最长的平衡子串 II
# https://leetcode.cn/problems/longest-balanced-substring-ii/description/
# 给你一个只包含字符 'a'、'b' 和 'c' 的字符串 s。
# 如果一个 子串 中所有 不同 字符出现的次数都 相同，则称该子串为 平衡 子串。
# 请返回 s 的 最长平衡子串 的 长度 。
# 子串 是字符串中连续的、非空 的字符序列。


from collections import defaultdict
from itertools import groupby


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def longestBalanced(self, s: str) -> int:
        res = 0

        for _, group in groupby(s):
            len_ = len(list(group))
            res = max2(res, len_)

        cands2 = [("a", "b"), ("a", "c"), ("b", "c")]
        for c1, c2 in cands2:
            pos = {0: -1}
            d = 0
            for i, c in enumerate(s):
                if c == c1:
                    d += 1
                elif c == c2:
                    d -= 1
                else:
                    d = 0
                    pos = {0: i}
                    continue
                if d in pos:
                    res = max2(res, i - pos[d])
                else:
                    pos[d] = i

        pos = {(0, 0): -1}
        counter = defaultdict(int)
        for i, c in enumerate(s):
            counter[c] += 1
            state = (counter["a"] - counter["b"], counter["a"] - counter["c"])
            if state in pos:
                res = max2(res, i - pos[state])
            else:
                pos[state] = i

        return res
