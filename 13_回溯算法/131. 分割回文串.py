# 给你一个字符串 s，请你将 s 分割成一些子串，使每个子串都是 回文串 。返回 s 所有可能的分割方案。
# 回文串 是正着读和反着读都一样的字符串。

# https://leetcode.cn/problems/palindrome-partitioning/
# 131. 分割回文串
# 1 <= s.length <= 16
# s 仅由小写英文字母组成

from typing import List


class Solution:
    def partition(self, s: str) -> List[List[str]]:
        def bt(index: int, path: List[str]) -> None:
            nonlocal res
            if index == n:
                res.append(path[:])
                return
            for j in range(index, n):
                sub = s[index : j + 1]
                if sub == sub[::-1]:
                    path.append(sub)
                    bt(j + 1, path)
                    path.pop()

        n = len(s)
        res = []
        bt(0, [])
        return res
