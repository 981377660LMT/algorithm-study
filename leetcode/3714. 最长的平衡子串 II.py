# 3714. 最长的平衡子串 II
# https://leetcode.cn/problems/longest-balanced-substring-ii/description/
# 给你一个只包含字符 'a'、'b' 和 'c' 的字符串 s。
# 如果一个 子串 中所有 不同 字符出现的次数都 相同，则称该子串为 平衡 子串。
# 请返回 s 的 最长平衡子串 的 长度 。
# 子串 是字符串中连续的、非空 的字符序列。
# !还有随机哈希做法

from typing import Tuple
from itertools import combinations


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def longestBalanced(self, s: str, charset="abc") -> int:
        def f(comb: Tuple[str, ...]) -> int:
            res = 0
            mp = {c: i for i, c in enumerate(comb)}
            counter = [0] * len(comb)
            first = {tuple(counter): -1}
            for i in range(len(s)):
                c = mp.get(s[i])
                if c is None:
                    counter = [0] * len(comb)
                    first = {tuple(counter): i}
                    continue
                counter[c] += 1
                min_ = min(counter)
                for j in range(len(counter)):
                    counter[j] -= min_
                key = tuple(counter)
                if key in first:
                    res = max2(res, i - first[key])
                else:
                    first[key] = i
            return res

        res = 0
        for size in range(1, len(charset) + 1):
            for comb in combinations(charset, size):
                res = max2(res, f(comb))
        return res
