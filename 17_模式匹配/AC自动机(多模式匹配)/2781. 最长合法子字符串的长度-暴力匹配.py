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
# !1 <= forbidden[i].length <= 10
# forbidden[i] 只包含小写英文字母。

# !双指针 O(sum(len(forbidden))+n*len(forbidden)^2)


from typing import List


class Solution:
    def longestValidSubstring(self, word: str, forbidden: List[str]) -> int:
        bad = set(forbidden)
        maxLen = max(map(len, forbidden))
        res, left, n = 0, 0, len(word)
        for right in range(n):
            # 新加入一个元素，尾部不能在bad里，否则这一段都无效
            for len_ in range(1, min(right - left + 1, maxLen) + 1):
                if word[right - len_ + 1 : right + 1] in bad:
                    left = right - len_ + 2
                    break
            res = max(res, right - left + 1)
        return res


# word = "leetcode", forbidden = ["de","le","e"]
assert Solution().longestValidSubstring(word="leetcode", forbidden=["de", "le", "e"]) == 4
